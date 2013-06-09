<?php
$memcache_obj = memcache_connect('localhost', 11212);

for($x = 0; $x < 10000; $x++) {

	memcache_set($memcache_obj, 'var_key' . $x, 'some variable' . $x, 0, 0);

	echo memcache_get($memcache_obj, 'var_key' . $x);
echo "\n";
}

?>
