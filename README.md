Overview
--------
This repository contains the BibTeX file and HTML templates that are used to
create the [Internet censorship bibliography](http://censorbib.nymity.ch).

Build it
--------

You first need [`bibliogra.py`](https://github.com/NullHypothesis/bibliograpy)
to turn the BibTeX file into an HTML bibliography.

Then, run the following commands to write the bibliography to `OUTPUT_DIR`.

    $ ./fetch_pdfs.py references.bib OUTPUT_DIR
    $ bibliogra.py -H header.tpl -F footer.tpl -f references.bib OUTPUT_DIR

Feedback
--------
Contact: Philipp Winter <phw@nymity.ch>  
OpenPGP fingerprint: `B369 E7A2 18FE CEAD EB96  8C73 CF70 89E3 D7FD C0D0`
