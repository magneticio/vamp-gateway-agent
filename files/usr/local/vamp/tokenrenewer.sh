#!/usr/bin/env bash

if [ -e /usr/local/vamp/token ]; then
  VAMP_KEY_VALUE_STORE_TOKEN="$( cat /usr/local/vamp/token )"
fi

TOKEN=${VAMP_KEY_VALUE_STORE_TOKEN}
URL=${VAMP_KEY_VALUE_STORE_CONNECTION}

echo "Token renewer started to check ${URL}"
while true; do
  SLEEP_DURATION=0
  CREATION_TTL=0
  TTL=0
  LOOKUP=$(curl -s --header "X-Vault-Token: ${TOKEN}"  ${URL}/v1/auth/token/lookup-self)
  CREATION_TTL=$(echo ${LOOKUP} | jq .data.creation_ttl)
  TTL=$(echo ${LOOKUP} | jq .data.ttl)
  echo "${TTL} seconds left for expiration, creation duration is ${CREATION_TTL} seconds"
  let "SLEEP_DURATION=${TTL}/2"
  if [ "${SLEEP_DURATION}" -lt 10 ]
  then
    echo "SLEEP_DURATION is set to 10"
    SLEEP_DURATION=10
  fi
  echo "{ \"increment\": \"${CREATION_TTL}\" }" > payload.json
  echo "Renewing the token"
  RESULT=$(curl -s --header "X-Vault-Token: ${TOKEN}"  --request POST --data @payload.json ${URL}/v1/auth/token/renew-self || CURL_RETURN_CODE=$?)
  echo "New lease_duration: $(echo $RESULT | jq .auth.lease_duration) seconds"
  echo "Wait for ${SLEEP_DURATION} seconds"
  sleep ${SLEEP_DURATION}
done
