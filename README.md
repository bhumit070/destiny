# Destiny

- Destiny is a file organization tool written in **Golang**. It automatically scans files in a given directory, creating folders based on their file extensions for efficient file management.

## Table of contents

- [Destiny](#destiny)
  - [Table of contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Features](#features)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Arguments](#arguments)
  - [Contributing](#contributing)

## Introduction

- Destiny simplifies file organization by scanning a specified directory and placing files into folders based on their extensions. Say goodbye to messy directories and welcome a more organized file structure.

## Features

- Scans files in a given directory.
- Automatically organizes files into folders based on their extensions.
- Creates folders on the fly if they don't exist.

## Installation

To install this cli tool run the below command in your terminal

- For **mac** and **linux** users

```bash
curl https://raw.githubusercontent.com/bhumit070/destiny/main/install.sh | bash
```

## Usage

```bash
destiny PATH_TO_DIRECTORY
```

- if no directory is passed it will use current directory

### Arguments

| Flag | Description                                                                 |
| :--- | :-------------------------------------------------------------------------- |
| -y   | If this flag is passed it won't ask for user confirmation before proceeding |
| -q   | it won't show how many files moved after command is done executing          |

## Contributing

- Fork the repository.
- Create a new branch for your feature or bug fix.
- Make your changes and submit a pull request.
