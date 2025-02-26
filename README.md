# ig2hugo - Instagram to Hugo Converter

A tool for fetching Instagram posts and converting them into Hugo-compatible Markdown files with Instagram shortcodes.

## Description

This tool uses the Instagram Graph API to retrieve all media from a specific Instagram account and converts them into Hugo-compatible Markdown files. Each Instagram post is saved as a separate Markdown file using the Instagram shortcode, allowing the content to be seamlessly displayed on a Hugo website.

## Requirements

- [Go](https://golang.org/dl/) version 1.16 or higher
- An Instagram Business or Creator account
- A Facebook Developer account
- Access to the Instagram Graph API (Access Token)

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/cubyverse/ig2hugo.git
cd ig2hugo

# Build
go build -o ig2hugo

# Optional: Copy to a PATH directory
# Linux/macOS:
sudo cp ig2hugo /usr/local/bin/
```

### Download Binary Files

You can also download the latest binary files from the [Releases page](https://github.com/cubyverse/ig2hugo/releases).

## Instagram API Configuration

To use this tool, you need an Instagram API Access Token and a User ID:

1. Create an app on [Meta for Developers](https://developers.facebook.com/)
2. Add the "Instagram Basic Display" product
3. Configure the app and obtain an Access Token
4. Find your Instagram User ID through the Graph API Explorer or the Developer Dashboard

## Usage

```bash
# Basic usage
./ig2hugo -token=YOUR_ACCESS_TOKEN -user=YOUR_USER_ID

# Specify custom output directory
./ig2hugo -token=YOUR_ACCESS_TOKEN -user=YOUR_USER_ID -output=path/to/output

# Display help
./ig2hugo -help
```

### Parameters

- `-token`: Instagram API Access Token (required)
- `-user`: Instagram User ID (required)
- `-output`: Output directory for Hugo Markdown files (default: "content/instagram")

## Output Format

The tool creates a Markdown file for each Instagram post in the output directory (default: `content/instagram/`). The filename follows the format `YYYY-MM-DD-POSTID.md`.

Each Markdown file contains:

1. Date and Tags
2. The Instagram shortcode with the correct Instagram post ID (extracted from the post's permalink)

Example:

```markdown
---
title: ""
date: 2023-05-12T14:30:45-07:00
draft: false
tags:
  - instagram
---

{{< instagram CgV8_TYuJZe >}}

```

> **Note:** The tool extracts the Instagram shortcode ID from the post's permalink (the part after `/p/` in the URL). This ensures compatibility with Hugo's Instagram shortcode.

## Integration with Hugo

1. Ensure your Hugo theme supports the Instagram shortcode or add a custom shortcode
2. Copy the contents of the output directory to your Hugo project

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details. 
