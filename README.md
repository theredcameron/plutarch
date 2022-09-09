# Plutarch
Track changes made to webpages.

## Prerequisites
To use this application, please ensure the following are installed on your computer.
* [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
* [Golang](https://go.dev/doc/install)  

## Instructions
Use the following instructions to run the program.
1. Clone this repository
2. Create a blank GitHub repository and make it visible to the public (ensure the computer running this application has permissions to run the `git push` command to this newly created blank repository)
3. Create a file named `sites-list` alongside the `main.go` file that will contain the list of websites to save. It should look like the following:
```
https://www.facebook.com
https://www.google.com
https://www.twitter.com
```
4. Run the `go build` command while in the same directory as the `main.go` file. This will create the executable used for this program
5. When running the executable, several parameters are required. They are explained here:
    * interval - The interval in **minutes**
    * author_email - The email address associated with the repository created in step 1.
    * author_name - The name of the author
    * journal_path - The path of the GitHub repository created in step 1. This is the same path you would use if you were cloning the empty repository. For example, if your empty template is cloned via SSH with the command `git clone git@github.com:testuser123/plutarchs-journal.git`, you would enter `git@github.com:testuser123/plutarchs-journal.git` as the journal path. See example below.
    ```
    ./plutarch -interval=3 -author_email="testuser123@whateveremail.com" -author_name="Test User" -journal_path="git@github.com:testuser123/plutarchs-journal.git"
    ```

## A Note on Logging
This program does not stop when retrieving a page's data fails, but it does log the failure. Since this is how it works, please be sure to check the logs, in the `logs` directory, to confirm that everything is in working order. Do this especially when you add, remove, or alter entries from the `sites-list` file.