simple command-line interface for jira
--------------------------------------

useful in its own right to avoid manually visiting jira website, or
for use in git hooks, etc...

for getting started by cloning, building, etc:

    git clone --recursive git@github.com:xoba/jira.git
    cd jira
    source goinit.sh
    ./install.sh

to get generic help, run the executable ```jira``` without any args.

to get help on listing, or actually listing your issues:

    jira list -help
    jira list -reporter mra
    
same story goes for creating issues, deleting issues, and commenting on them.
    
you should put exports like these into your ~/.bashrc file for convenience:

    export JIRA_USERNAME=joe.smith
    export JIRA_PASSWORD=abc123
    export JIRA_URL=http://www.example.com/jira
