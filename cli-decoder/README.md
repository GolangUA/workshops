
## Test task:
Develop a cli tool that will:

- that parses cmd line tools. It should be either `json` or `xml`. 
- if JSON is selected - cmd line tool should only accept valid json string, same goes to xml.
- After data is received and validated it should be stored as file in current dir with json or xml format.
- Data is saved only if it's unique. 
- Uniquness should be identified using md5 has of a file.
- On startup identify all known md5 hashsums
- Put hash data into the separate file and keep them in memory when program is running  
- If error appeared - program should still be working and accept new data. Handle all the errors and panics properly

## Preparation steps
0. Initialize go mod
1. Install https://github.com/spf13/cobra or https://github.com/urfave/cli and read the documentation how to worl with these cli-tools libs.
2. Create main and write code to parse flags.
3. Implement the interaction with files and md5
4. Implement json and xml validators
5. Run your solution:
`go run main.go -json`
6. Use json to validate if it works   
`'{"alias":"go-dms-workshop","desc":"Create app and try it with different DMS", "type":"important", "ts":1473837996,"tags":["Golang","Workshop","DMS"],"etime":"4h","rtime":"8h","reminders":["3h", "15m"]}'`
7. Validate with xml
