<template>
  <div class="mt-6 flex w-full justify-center">
    <div
      class="w-1/2 justify-center rounded-xl border-t-2 border-entrust-purple bg-white shadow-lg"
    >
      <div class="flex w-full justify-center border-b-2 border-gray-300">
        <LogosIDaaS class="max-w-3/4 mx-6 h-40" />
      </div>
      <slot />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import Web3 from "web3";
const config = useRuntimeConfig();
const idpBackUrl = config.public.IDP_SERVER_URL;
// var nonce = getCookie("nonce");
if (process.client) {
  var nonce = window.sessionStorage.getItem("nonce");
}
var accounts;

// connectWallet();

// //state
// const wallet = ref(" ");

async function connectWallet() {
  if (typeof window !== "undefined") {
    if (window.ethereum) {
      await ethereum.request({
        method: "wallet_requestPermissions",
        params: [{ eth_accounts: {} }],
      });
      accounts = await window.ethereum.request({
        method: "eth_requestAccounts",
      });
      const walletAddress = accounts[0];
      //document.getElementById("walletId").innerHTML = `Wallet connected: ${walletAddress}`
      msgParams.message.signer = walletAddress;
      console.log(walletAddress);
      wallet.value = "Wallet connected: " + walletAddress;
    } else {
      wallet.value = "Wallet not found";
    }
  }
}

function onSignClick() {
  console.log("CLICK");
  var params = [accounts[0], JSON.stringify(msgParams)];
  var method = "eth_signTypedData_v4";
  web3.currentProvider.sendAsync(
    {
      method,
      params,
      from: accounts[0],
    },
    function (err, result) {
      if (err) return console.dir(err);
      if (result.error) {
        alert(result.error.message);
      }
      if (result.error) return console.error("ERROR", result);
      console.log("RESULT:" + JSON.stringify(result));
      let signatureJSON = {
        message: msgParams,
        signature: result.result,
      };
      console.log(`SIGNATURE JSON: ${JSON.stringify(signatureJSON)}`);
      fetch(idpBackUrl + "/api/oidc/enforce", {
        method: "POST",
        mode: "cors",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(signatureJSON),
      })
        .then((response) => {
          if (response.ok) {
            window.location.href = response.url;
          } else {
            console.log("Error: ", response.status);
          }
        })
        .catch((error) => {
          console.error(error);
        });
    },
  );
}

function getCookie(name) {
  // Split cookie string and get all individual name=value pairs in an array
  if (typeof window !== "undefined") {
    var cookieArr = document.cookie.split(";");

    // Loop through the array elements
    for (var i = 0; i < cookieArr.length; i++) {
      var cookiePair = cookieArr[i].split("=");

      /* Removing whitespace at the beginning of the cookie name and compare it with the given string */
      if (name == cookiePair[0].trim()) {
        // Decode the cookie value and return
        return decodeURIComponent(cookiePair[1]);
      }
    }

    // Return null if not found
    return null;
  }
}

function getTimeAndDate() {
  let today = new Date();
  let date =
    today.getFullYear() + "-" + (today.getMonth() + 1) + "-" + today.getDate();
  let time =
    today.getHours() + ":" + today.getMinutes() + ":" + today.getSeconds();
  let dateTime = date + " " + time;
  return dateTime;
}
const msgParams = {
  domain: {
    chainId: BigInt(1337).toString(),
    name: "NFT ownership verification service",
    verifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
    version: "1",
  },
  message: {
    why: "To prove you have a NFT of the mentioned above collection, please sign this message",
    nft_collection: "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D",
    signer: "",
    timestamp: getTimeAndDate(),
    nonce: nonce,
  },
  // Refers to the keys of the *types* object below.
  primaryType: "Message to sign",
  types: {
    // Clarify if EIP712Domain refers to the domain the contract is hosted on
    EIP712Domain: [
      { name: "name", type: "string" },
      { name: "version", type: "string" },
      { name: "chainId", type: "uint256" },
      { name: "verifyingContract", type: "address" },
    ],

    "Message to sign": [
      { name: "why", type: "string" },
      { name: "nft_collection", type: "address" },
      { name: "signer", type: "address" },
      { name: "timestamp", type: "string" },
      { name: "nonce", type: "string" },
      //{name: 'exp', type: 'string'}
    ],
  },
};
</script>
