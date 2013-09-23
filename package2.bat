matsano p2 pad "YELLOW SUBMARINE"
matsano p2 decrypt "YELLOW SUBMARINE" https://gist.github.com/tqbf/3132976/raw/f0802a5bc9ffa2a69cd92c981438399d4ce1b8e4/gistfile1.txt
matsano p2 randomEncrypt file:package1\gistfile4.txt | matsano p2 blockMode
matsano p2 randomEncrypt file:package1\gistfile4.txt | matsano p2 blockMode
matsano p2 randomEncrypt file:package1\gistfile4.txt | matsano p2 blockMode
matsano p2 crackAesEcb "Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK"
matsano p2 randomKey > key.txt
matsano p2 profileCrack file:key.txt admin | matsano p2 profileDecrypt file:key.txt