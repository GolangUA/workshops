
How to read data from std input?
```
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	// when you hit enter - the data is read
	scanner.Scan()
	// this will read data from the scanner
	yourData := scanner.Text()
	// will print it to stdin
	fmt.Println(yourData)
}

```

Try to use it in your tasks

****Task 0****

Print a chessboard with the specified dimensions of height and width, according to the example height - 4 and wigth 6:

```
^  ^  ^  ^  ^  ^
  ^  ^  ^  ^  ^  ^
^  ^  ^  ^  ^  ^
  ^  ^  ^  ^  ^  ^
```

Input parameters: height, width, symbol to print
Output: board

****Task 1****


scan array string

use strings split to split numbers and Atoi

add logic that counts the number of positive even numbers in the array and prints it.

        Example of input data: 30,-1,-6,90,-6,22,52,123,2,35,6

****Task 2****

1. Enter valid bank card number
2. Validate it
3. Print string with all strings covered except for the last 4

        Input example: 4539 1488 0343 6467


****Task 3****


Implement a fibonacci function that returns a function (a closure) that returns successive fibonacci numbers (0, 1, 1, 2, 3, 5, …)
The next number is found by adding up the two numbers before it:

the 2 is found by adding the two numbers before it (1+1),
the 3 is found by adding the two numbers before it (1+2),
the 5 is (2+3)


****Task 4****

Check if a given number or part of this number is a palindrome. For example, the number 1234437 is not a palindrome,
but its part 3443 is a palindrome. Numbers less than 10 are counted as invalid input.
```
Easy level:
Find if a number is a palindrome or not

Mid level:
You need to find only first subpalindrome

Hard:
Find all subpalindromes
```

Input parameters: number

Output: the palindrome/s extracted from the number, or 0 if the extraction failed or no palindromes was found.


****Task 5****

There are 2 ways to count lucky tickets:
1. Simple - if a six-digit number is printed on the ticket, and the sum of the first three digits is equal to the sum of the last three digits, then this
   ticket is lucky.
2. Hard - if the sum of the even digits of the ticket is equal to the sum of the odd digits of the ticket, then the ticket is considered lucky.
   Determine programmatically which variant of counting lucky tickets will give them a greater number at a given interval.

Task: Calculate how many tickets are lucky within provided min and max tickets

Input parameters: 2 values min digit of ticket and max digit of ticket
Output: information about the winning method and the number of lucky tickets for each counting method.

Input numbers contains exactly 6 digits. Not more or less


Example:
```
Min: 120123
Max: 320320

--Result--
EasyFormula: 11187
HardFormula: 5790
```
