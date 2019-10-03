fork {
	puts "child: #{Process.pid}"
	sleep
}

puts "parent: #{Process.pid}"
sleep
