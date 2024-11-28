# HashCracker

## App context:

The hash cracker is a cmd tool that will read a file containg passwords, and hashed password

The cracker must “crack” the password by using it’s passwords list

## App flow:

1. Read word list (may need to be done in chunks): https://github.com/danielmiessler/SecLists/blob/master/Passwords/Leaked-Databases/rockyou-75.txt
2. Hash every word in the file and try to match it with the provided hash
3. If a match is found, the app prints the password  and exits

## App requirements:

1. Everything must be handled as concurrently as possible
2. The supported hashing algorithms are: SHA256, keccak256, MD5
3. Think of a cool name for the app!
4. Add unit tests

# How to set up the project

## 1. Clone the Project

Clone the project repository to your local machine:

```bash
git clone https://github.com/Dosik13/PasswordDestroyer.git
```

## 2. Start the program

```bash
cracker <path_to_wordlist> <hash_to_crack> --debug
```

`debug` is an optional flag and enables debug logs

# Example usage:

Input:

```bash
cracker passwords.txt 3fc0a7acf087f549ac2b266baf94b8b1
```

Output:
"Match found!","password":"qwerty123"
