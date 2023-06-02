# assignment_demo_2023

![Tests](https://github.com/TikTokTechImmersion/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

This is the backend assignment for the 2023 TikTok Tech Immersion Program.

## Setup Instructions

1. Clone the repository [here](https://github.com/melvincwng/tiktok_assignment_demo_2023) using this command:

   > `git clone https://github.com/melvincwng/tiktok_assignment_demo_2023.git --config core.autocrlf=false`.

2. The flag --config core.autocrlf=false is to ensure that the line endings are in Unix format (LF) and not Windows format (CRLF). This is to ensure that the bash scripts can be executed properly. If this is not done, you may encounter issues regarding `/usr/bin/env: ‘bash\r’: No such file or directory` when running `docker-compose up` later on.

3. Download [Go](https://go.dev/doc/install).

4. Download [Docker Desktop](https://www.docker.com/products/docker-desktop/).

5. Run the command `docker-compose up` in the root directory of the project.

6. For Windows Users:

   - After executing `docker-compose up`, you may encounter issues regarding `build.sh not found` despite the file being present in the directory.
   - This is because the line endings **are not** in Unix format (LF).
   - Bash scripts are **sensitive** to line endings formatting (can read up on LF vs CRLF).
   - To fix this, on your code editor e.g. VSCode, `click on the CRLF button at the bottom right corner and change it to LF`.
   - Then, run `docker-compose up` again. It should work now. 
   - You can change your settings back to CRLF after you are done.
   - See **References 1 and 2 below** for more information.

7. Go to Docker Desktop and you should now see your containers running over there.

## References

How to resolve CRLF/LF line endings related issues for Windows users:

1. https://willi.am/blog/2016/08/11/docker-for-windows-dealing-with-windows-line-endings/
2. https://stackoverflow.com/questions/39527571/are-shell-scripts-sensitive-to-encoding-and-line-endings
