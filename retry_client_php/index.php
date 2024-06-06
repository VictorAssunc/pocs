<?php

include("Client.php");

$client = new Client("localhost:4000");
$client->sendRequest("/test/42?a=0");
$client->sendRequest("/test?a=0&b=42");
