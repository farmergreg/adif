#!/usr/bin/env sh

curl -v --location 'localhost:8080/api/v1/' \
--header 'Content-Type: text/x-adif-adi' \
--data '<ADIF_VER:5>3.1.5<EOH>

<CALL:5>K9CTS<BAND:3>40m<MODE:2>CW<EOR>
<CALL:5>W9PVA<BAND:3>40m<EOR>'