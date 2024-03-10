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
  <link rel="icon" type="image/svg+xml" href="favicon.svg">
  <style>
  body {
    font-family: Roboto,Helvetica,sans-serif;
    background: #ddd;
    margin-left: auto;
    margin-right: auto;
    margin-top: 20px;
    max-width: 1000px;
  }
  li {
    margin-top: 1em;
    margin-bottom: 1em;
    margin-right: 1em;
  }
  h1 {
    font-size: 25px;
    color: #efefef;
    width: 80%;
    float: left;
  }
  ul {
    border-radius: 10px;
    border:1px solid #c0c0c0;
    background: #efefef;
    box-shadow: 2px 2px 5px #bbb;
  }
  ul.a {
    list-style-image: url('donate-icon.svg');
  }
  a:link {
    color:#0b61a4;
    text-decoration:none;
  }
  a:visited {
    color:#033e6b;
    text-decoration:none;
  }
  a:hover {
    text-decoration:underline;
  }
  p {
    margin: 0px;
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
    margin-right: 0.5em;
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
    background-color: #f6ba81;
  }
  #left-header {
    flex: 4;
    background: #efefef;
    margin-right: 0.5em;
    border-radius: 10px;
    border: 1px solid #c0c0c0;
    box-shadow: 2px 2px 5px #bbb;
    overflow: hidden; /* For child elements to inherit rounded corners. */
  }
  #right-header {
    flex: 1;
    background: #efefef;
    margin-left: 0.5em;
    background: #333 url('img/research-power-tools-cover.jpg') no-repeat;
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
    background: #333 url('open-access.svg') right/25% no-repeat;
  }
  #censorbib-description {
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
    color: #d94b7b
  }
  #book-info > a:visited {
    color: #d94b7b
  }
  </style>
</head>

<body>

  <div class="flex-row">

    <div id="left-header" class="flex-column round-shadow">

      <div id="title-box">
        <h1>Selected Research Papers<br>in Internet Censorship</h1>
      </div>

      <div class="flex-row">

        <div id="censorbib-description">
          CensorBib is an online archive of selected research papers in the field
          of Internet censorship.  Most papers on CensorBib approach the topic
          from a technical angle, by proposing designs that circumvent censorship
          systems, or by measuring how censorship works.  The icons next to each
          paper make it easy to download, cite, and link to papers.  If you think
          I missed a paper,
          <a href="https://github.com/NullHypothesis/censorbib">
          make a pull request
          </a>.
          Finally, the
          <a href="https://github.com/net4people/bbs/issues">net4people/bbs forum</a>
          has reading groups for many of the papers listed below.
        </div> <!-- censorbib-description -->

        <div id="censorbib-links">
          <div class="menu-item">
              <img class="top-icon" src="img/lock-icon.svg" alt="onion service icon">
              <a href="http://putnst3yv7k6vvb3avdqgdutrz3kaufitaiwbjhjox7o3daakr43fhad.onion">Onion service mirror</a>
          </div>
          <div class="menu-item">
            <img class="top-icon" src="img/code-icon.svg" alt="source code icon">
            <a href="https://github.com/NullHypothesis/censorbib">CensorBib code</a>
          </div>
          <div class="menu-item">
            <img class="top-icon" src="img/update-icon.svg" alt="update icon">
            <a href="https://github.com/NullHypothesis/censorbib/commits/master">Last update: {{.Date}}</a>
          </div>
          <div class="menu-item">
            <img class="top-icon" src="img/donate-icon.svg" alt="donate icon">
            <a href="https://nymity.ch/donate.html">Donate</a>
          </div>
        </div> <!-- censorbib-links -->

      </div>

    </div> <!-- left-header -->

    <div id="right-header" class="round-shadow">

      <div class="flex-column" style="height: 100%">
        <div style="flex: 1 1 auto">
        </div>

        <div id="book-info" style="flex: 0 1 auto">
          Are you a researcher? If so, you may like my book
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
