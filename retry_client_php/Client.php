<?php

/**
 * Class Client.
 */
class Client {
    /**
     * The default timeout, in milliseconds, for each request.
     */
    const DEFAULT_TIMEOUT = 30000;

    /**
     * The map of custom timeouts, in milliseconds, for endpoints.
     */
    const TIMEOUT_MAP = array(
        "/test\/[0-9]+(?:\?.*)?$/" => 1100,
    );

    /**
     * The maximum number of retries for each request.
     */
    const MAX_RETRIES = 8;

    /**
     * The minimum interval, in microseconds, between retries.
     */
    const MIN_RETRY_INTERVAL = 10000;

    /**
     * The maximum interval, in microseconds, between retries.
     */
    const MAX_RETRY_INTERVAL = 100000;

    /**
     * The randomization factor for jitter the retry interval.
     */
    const RANDOMIZATION_FACTOR = 0.3;

    /**
     * @var string
     */
    private $basePath;

    public function __construct($basePath = '') {
        if ($basePath != '') {
            $this->basePath = $basePath;
        }
    }

    /**
     * Envia uma requisição para a URL base concatenada com $path.
     *
     * @param string   $path
     * @param string   $jwt
     * @param array    $data
     * @param string   $type
     * @param string[] $extraHeaders
     *
     * @return array
     */
    public function sendRequest($path, $jwt = '', $data = null, $type = 'GET', $extraHeaders = []) {
        $headers = [
            'Content-Type: application/json;charset=UTF-8',
            'X-Authorization:' . $jwt,
        ];

        if (!empty($extraHeaders)) {
            $headers = array_unique(array_merge($headers, $extraHeaders));
        }

        $bpc = preg_match('/\?/i', $path) ? '&' : '?';
        $bpc = $bpc . 'bpc=true';

        $options = [
            CURLOPT_URL => $this->basePath . $path . $bpc,
            CURLOPT_HTTPHEADER => $headers,
            CURLOPT_RETURNTRANSFER => 1,
            CURLOPT_TIMEOUT_MS => self::DEFAULT_TIMEOUT,
        ];

        $retry = false;
        if (!empty($data) && $type == 'POST') {
            $options[CURLOPT_POST] = 1;
            $options[CURLOPT_POSTFIELDS] = json_encode($data);
        } else if ($type != 'POST') {
            $matchingPattern = array_filter(array_keys(self::TIMEOUT_MAP), function($pattern) use ($path) {
                return preg_match($pattern, $path);
            });
            if (!empty($matchingPattern)) {
                $retry = true;
                $options[CURLOPT_TIMEOUT_MS] = self::TIMEOUT_MAP[reset($matchingPattern)];
            }
        }

        echo $options[CURLOPT_TIMEOUT_MS];
        $response = '';
        $info = [];
        $ok = false;
        $timeoutError = $httpCode = $curlErr = 0;
        for ($i = 0; $i < self::MAX_RETRIES && !$ok; $i++) {
            $ch = curl_init();
            curl_setopt_array($ch, $options);

            $response = trim(curl_exec($ch));
            $info = curl_getinfo($ch);
            $httpCode = $info["http_code"];

            $curlErr = curl_errno($ch);
            curl_close($ch);

            if ($curlErr != 0 || $httpCode >= 500) {
                if ($curlErr == CURLE_OPERATION_TIMEDOUT) {
                    $timeoutError++;
                }

                if ($retry) {
                    $delay = self::MIN_RETRY_INTERVAL * (2 ** $i);
                    if ($delay > self::MAX_RETRY_INTERVAL) {
                        $delay = self::MAX_RETRY_INTERVAL;
                    }

                    $delta = self::RANDOMIZATION_FACTOR * $delay;
                    usleep(rand($delay - $delta, $delay + $delta));
                    continue;
                }
            }

            $ok = true;
        }

        if ($i > 0 && $timeoutError > 0) {
            trigger_error("Attempts for '$path': $i - Timeout errors: $timeoutError - Last Status: $httpCode");
        }

        if ($httpCode == 0 || ($httpCode >= 400 && $httpCode <= 599)) {
            trigger_error("Erro ao realizar a requisição: Erro: Status code $httpCode - Curl error: $curlErr");
            $httpCode = 500;
        }

        echo "$response\n";
        return [
            'info' => $info,
            'response' => json_decode($response, true),
        ];
    }
}