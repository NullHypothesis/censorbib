package main

import (
	"bytes"
	"html/template"
	"log"
	"time"
)

const headerTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <meta name="color-scheme" content="light dark">
  <title>The Internet censorship bibliography</title>
  <link rel="icon" href="assets/favicon-32.png"  sizes="32x32">
  <link rel="icon" href="assets/favicon-128.png" sizes="128x128">
  <link rel="icon" href="assets/favicon-180.png" sizes="180x180">
  <link rel="icon" href="assets/favicon-192.png" sizes="192x192">
  <style>
  :root {
    color-scheme: light dark;
    --page: #ddd;
    --surface: #f5f5f5;
    --text: #222;
    --muted: #666;
    --border: #c0c0c0;
    --shadow: #bbb;
    --link: #0b61a4;
    --visited: #033e6b;
    --highlight: #ffb772;
    --mark-text: #222;
    --input: #fff;
    --input-border: #aaa;
  }
  @media (prefers-color-scheme: dark) {
    :root {
      --page: #171c22;
      --surface: #222a33;
      --text: #e7e5df;
      --muted: #b2bac4;
      --border: #404c5a;
      --shadow: #101419;
      --link: #91c5ed;
      --visited: #c4b2df;
      --highlight: #493a28;
      --mark-text: #f5dfb8;
      --input: #1b222a;
      --input-border: #596779;
    }
    .icon,
    .top-icon {
      filter: invert(88%);
    }
    .icons .icon:hover {
      background-color: transparent;
    }
    .icons a:hover {
      background-color: var(--highlight);
      border-radius: 25%;
    }
  }
  body {
    font-family: Roboto, Helvetica, sans-serif;
    background: var(--page);
    color: var(--text);
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
    border: 1px solid var(--border);
    background: var(--surface);
    box-shadow: 2px 2px 5px var(--shadow);
  }
  a:link {
    color: var(--link);
    text-decoration: none;
  }
  a:visited {
    color: var(--visited);
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
    color: var(--muted);
  }
  .venue {
    font-style: italic;
  }
  .paper-entry {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    column-gap: 0.75em;
  }
  .paper-entry > :not(.icons) {
    grid-column: 1;
  }
  .paper-entry[hidden] {
    display: none;
  }
  .paper {
    min-width: 0;
    overflow-wrap: anywhere;
    font-weight: bold;
  }
  .other {
    color: var(--muted);
  }
  #container mark {
    color: var(--mark-text);
    background: var(--highlight);
    border-radius: 2px;
    padding: 0 0.1em;
  }
  #footer {
    text-align: center;
    line-height: 20px;
  }
  .icon {
    height: 1em;
    width: 1em;
    display: block;
    border-radius: 25%;
    transition: background-color 0.3s ease;
    padding: 0.5em;
    overflow: visible;
  }
  .icon:hover {
    background-color: var(--highlight);
    cursor: pointer;
  }
  .icons a {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    line-height: 0;
  }
  a:focus-visible,
  input:focus-visible,
  button:focus-visible {
    outline: 2px solid var(--link);
    outline-offset: 2px;
  }
  .icons {
    grid-column: 2;
    grid-row: 1 / span 3;
    align-self: start;
    display: inline-flex;
    align-items: center;
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
    background-color: var(--highlight);
  }
  #header {
    margin: 1em;
  }
  #left-header {
    flex: 4;
    background: var(--surface);
    margin-right: 0.5em;
    border-radius: 10px;
    border: 1px solid var(--border);
    box-shadow: 2px 2px 5px var(--shadow);
    overflow: hidden; /* For child elements to inherit rounded corners. */
  }
  #right-header {
    flex: 1;
    background: var(--surface);
    margin-left: 0.5em;
    background: #333 url('assets/research-power-tools-cover.jpg') no-repeat;
    background-size: 100%;
  }
  .round-shadow {
    border-radius: 10px;
    border: 1px solid var(--border);
    box-shadow: 2px 2px 5px var(--shadow);
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
  #search-form {
    display: flex;
    align-items: center;
    gap: 0.75em;
    margin: 1em;
    padding: 0.75em;
    background: var(--surface);
    border-radius: 10px;
    border: 1px solid var(--border);
    box-shadow: 2px 2px 5px var(--shadow);
  }
  #search-form label {
    font-weight: bold;
  }
  #search-input {
    flex: 1;
    min-width: 0;
    font: inherit;
    padding: 0.5em 0.75em;
    border: 1px solid var(--input-border);
    border-radius: 6px;
    background: var(--input);
    color: var(--text);
  }
  #result-count {
    color: var(--muted);
    white-space: nowrap;
  }
  #no-results {
    margin: 1em;
    padding: 1em;
    text-align: center;
    color: var(--muted);
    background: var(--surface);
    border-radius: 10px;
    border: 1px solid var(--border);
    box-shadow: 2px 2px 5px var(--shadow);
  }
  #bibtex-modal[hidden] {
    display: none;
  }
  #bibtex-modal {
    position: fixed;
    inset: 0;
    z-index: 100;
  }
  #bibtex-backdrop {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.55);
  }
  #bibtex-dialog {
    position: relative;
    max-width: min(760px, calc(100vw - 2em));
    max-height: calc(100vh - 2em);
    margin: 1em auto;
    display: flex;
    flex-direction: column;
    background: var(--surface);
    border-radius: 10px;
    border: 1px solid var(--border);
    box-shadow: 2px 2px 10px #333;
  }
  #bibtex-dialog header {
    display: flex;
    align-items: center;
    gap: 0.5em;
    padding: 0.75em;
    border-bottom: 1px solid var(--border);
  }
  #bibtex-title {
    flex: 1;
    margin: 0;
    color: var(--text);
    font-size: 1.1em;
  }
  #bibtex-copy-status {
    color: var(--muted);
    min-width: 4.5em;
    text-align: right;
  }
  #bibtex-content {
    margin: 0;
    padding: 1em;
    overflow: auto;
    white-space: pre-wrap;
    font-size: 0.9em;
    line-height: 1.35;
    background: var(--input);
    color: var(--text);
  }
  @media (prefers-color-scheme: dark) {
    #book-info > a:link,
    #book-info > a:visited {
      color: #efa5bd;
    }
  }
  @media (max-width: 720px) {
    .paper-entry {
      grid-template-columns: minmax(0, 1fr);
    }
    .icons {
      grid-column: 1;
      grid-row: 4;
      margin-top: 0.25em;
      flex-wrap: wrap;
      max-width: 100%;
    }
    #header,
    .flex-row {
      flex-direction: column;
    }
    #left-header,
    #right-header {
      margin: 0 0 1em 0;
    }
    #search-form {
      align-items: stretch;
      flex-direction: column;
    }
    #result-count {
      white-space: normal;
    }
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

var headerTmpl = template.Must(template.New("header").Parse(headerTemplate))

func header() string {
	i := struct {
		Date string
	}{
		Date: time.Now().UTC().Format(time.DateOnly),
	}
	buf := new(bytes.Buffer)
	if err := headerTmpl.Execute(buf, i); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
	return buf.String()
}
