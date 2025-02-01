package main

import (
	"bytes"
	"log"
	"text/template"
	"time"
)

const headerTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <title>The Internet censorship bibliography</title>
  <link rel="icon" href="assets/favicon-32.png"  sizes="32x32">
  <link rel="icon" href="assets/favicon-128.png" sizes="128x128">
  <link rel="icon" href="assets/favicon-180.png" sizes="180x180">
  <link rel="icon" href="assets/favicon-192.png" sizes="192x192">
  <style>
  body {
    font-family: Roboto, Helvetica, sans-serif;
    background: #ddd;
    margin: 1em auto;
    max-width: 1000px;
  }
  li {
    border-radius: 10px;
    margin: 0.5em;
    padding: 0.5em;
  }
  h1 {
    font-size: 2em;
    color: #efefef;
    width: 80%;
    float: left;
  }
  ul {
    padding: 0.5em;
    list-style-type: none; /* Disable bullet points */
    border-radius: 10px;
    border: 1px solid #c0c0c0;
    background: #f5f5f5;
    box-shadow: 2px 2px 5px #bbb;
  }
  a:link {
    color: #0b61a4;
    text-decoration: none;
  }
  a:visited {
    color: #033e6b;
    text-decoration: none;
  }
  a:hover {
    text-decoration: underline;
  }
  .icons a:hover {
    text-decoration: none;
  }
  p {
    margin: 0px;
  }
  #container {
    margin: 1em;
  }
  .author {
    color: #666;
  }
  .venue {
    font-style: italic;
  }
  .paper {
    font-weight: bold;
  }
  .other {
    color: #666;
  }
  #footer {
    text-align: center;
    line-height: 20px;
  }
  .icon {
    height: 1em;
    width: 1em;
    display: inline-block;
    border-radius: 50%;
    transition: background-color 0.3s ease;
    padding: 0.5em;
    overflow: visible;
  }
  .icon:hover {
    background-color: #ffb772;
    cursor: pointer;
  }
  .icons {
    float: right;
  }
  .top-icon {
    height: 1em;
    width: 1em;
    position: relative;
    vertical-align: middle;
    margin-left: 1em;
  }
  .menu-item {
    padding-bottom: 5px;
  }
  .url {
    font-family: monospace;
    font-size: 12px;
  }
  :target {
    background-color: #ffb772;
  }
  #header {
    margin: 1em;
  }
  #left-header {
    flex: 4;
    background: #f5f5f5;
    margin-right: 0.5em;
    border-radius: 10px;
    border: 1px solid #c0c0c0;
    box-shadow: 2px 2px 5px #bbb;
    overflow: hidden; /* For child elements to inherit rounded corners. */
  }
  #right-header {
    flex: 1;
    background: #f5f5f5;
    margin-left: 0.5em;
    background: #333 url('assets/research-power-tools-cover.jpg') no-repeat;
    background-size: 100%;
  }
  .round-shadow {
    border-radius: 10px;
    border: 1px solid #c0c0c0;
    box-shadow: 2px 2px 5px #bbb;
    overflow: hidden; /* For child elements to inherit rounded corners. */
  }
  .flex-row {
    display: flex;
  }
  .flex-column {
    display: flex;
    flex-direction: column;
  }
  #title-box {
    text-align: center;
    background: #333 url('assets/open-access.svg') right/25% no-repeat;
  }
  #censorbib-description {
    font-size: 1.15em;
    text-align: justify;
    padding: 1em;
    flex: 5;
  }
  #censorbib-links {
    padding: 1em;
    flex: 2;
    font-size: 0.9em;
  }
  #book-info {
    text-align: center;
    padding: 0.5em;
    background: #333;
    color: #efefef;
  }
  #book-info > a:link {
    color: #d94b7b;
  }
  #book-info > a:visited {
    color: #d94b7b;
  }
  </style>
</head>

<body>

  <div id="header" class="flex-row">

    <div id="left-header" class="flex-column round-shadow">

      <div id="title-box">
        <h1>Selected Research Papers<br>in Internet Censorship</h1>
      </div>

      <div class="flex-row">

        <div id="censorbib-description">
          CensorBib is an archive of selected academic research papers on
          Internet censorship.  If you think I missed a paper,
          <a href="https://github.com/NullHypothesis/censorbib">make a pull request</a>.
          Finally, the
          <a href="https://github.com/net4people/bbs/issues">net4people/bbs forum</a>
          has reading groups for many of the papers listed below.
        </div> <!-- censorbib-description -->

        <div id="censorbib-links">
          <div class="menu-item">
            <img class="top-icon" src="assets/code-icon.svg" alt="source code icon">
            <a href="https://github.com/NullHypothesis/censorbib">CensorBib code</a>
          </div>
          <div class="menu-item">
            <img class="top-icon" src="assets/update-icon.svg" alt="update icon">
            <a href="https://github.com/NullHypothesis/censorbib/commits/master">Updated: {{.Date}}</a>
          </div>
        </div> <!-- censorbib-links -->

      </div>

    </div> <!-- left-header -->

    <div id="right-header" class="round-shadow">

      <div class="flex-column" style="height: 100%">
        <div style="flex: 1 1 auto">
        </div>

        <div id="book-info" style="flex: 0 1 auto">
          Are you a researcher? You may like my book
          <a href="http://research-power-tools.com">Research Power Tools</a>.
        </div>
      </div>

    </div> <!-- right-header -->

  </div>`

func header() string {
	tmpl, err := template.New("header").Parse(headerTemplate)
	if err != nil {
		log.Fatal(err)
	}
	i := struct {
		Date string
	}{
		Date: time.Now().UTC().Format(time.DateOnly),
	}
	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, i); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
	return buf.String()
}
