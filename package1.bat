matsano p1 hex base64 49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d
matsano p1 base64 hex SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t
matsano p1 fixedXor 1c0111001f010100061a024b53535009181c 686974207468652062756c6c277320657965
matsano p1 decryptSingleXor 1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736
matsano p1 detectSingleXor https://gist.github.com/tqbf/3132713/raw/gistfile1.txt
(echo Burning 'em, if you ain't quick and nimble && echo I go crazy when I hear a cymbal) | matsano p1 xor ICE
(echo Burning 'em, if you ain't quick and nimble && echo I go crazy when I hear a cymbal) | matsano p1 xor ICE | matsano p1 -ei hex -eo ascii xor ICE
matsano p1 decryptXor https://gist.github.com/tqbf/3132752/raw/gistfile1.txt
matsano p1 decryptAes "YELLOW SUBMARINE" https://gist.github.com/tqbf/3132853/raw/gistfile1.txt
matsano p1 detectAes https://gist.github.com/tqbf/3132928/raw/gistfile1.txt

