#!/usr/bin/env python3
#
# Copyright 2015 Philipp Winter <phw@nymity.ch>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
"""
Fetch pdf and ps files in BibTeX file.
"""

import os
import sys
import errno
import urllib.request

import pybtex.database.input.bibtex as bibtex


def download_pdf(url, file_name):
    """
    Download file and write it to given file name.
    """

    print("Now fetching %s" % url)

    try:
        fetched_file = urllib.request.urlopen(url)
    except Exception as err:
        print(err, file=sys.stderr)
        return

    with open(file_name, "wb") as fd:
        fd.write(fetched_file.read())


def main(file_name, output_dir):
    """
    Extract BibTeX key and URL, and then trigger file download.
    """

    parser = bibtex.Parser()
    bibdata = parser.parse_file(file_name)

    # Create download directories.

    try:
        os.makedirs(os.path.join(output_dir, "pdf"))
        os.makedirs(os.path.join(output_dir, "ps"))
    except OSError as exc:
        if exc.errno == errno.EEXIST:
            pass
        else:
            raise

    # Iterate over all BibTeX entries and trigger download if necessary.

    for bibkey in bibdata.entries:

        entry = bibdata.entries[bibkey]
        url = entry.fields.get("url")
        if url is None:
            continue

        # Extract file name extension and see what we are dealing with.

        _, ext = os.path.splitext(url)
        if ext:
            ext = ext[1:]

        if ext not in ["pdf", "ps"]:
            print("Skipping %s because it's not a pdf or ps file." % url,
                  file=sys.stderr)
            continue

        file_name = os.path.join(output_dir, ext, bibkey + ".%s" % ext)
        if os.path.exists(file_name):
            print("Skipping %s because we already have it." % file_name,
                  file=sys.stderr)
            continue

        download_pdf(url, file_name)

    return 0


if __name__ == "__main__":

    if len(sys.argv) != 3:
        print("\nUsage: %s FILE_NAME OUTPUT_DIR\n" % sys.argv[0],
              file=sys.stderr)
        sys.exit(1)

    sys.exit(main(sys.argv[1], sys.argv[2]))
