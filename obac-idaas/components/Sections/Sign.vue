<template>
  <div class="mt-8 w-full px-8" v-if="!isVerifyingSignature">
    <h2 class="ml-2 text-3xl font-bold text-gray-700">
      Hi, fifa.com wants you to prove ownership of an asset
    </h2>
    <h2 class="ml-2 mt-3 text-xl font-bold">
      Access is granted to owners of the following smart contracts:
    </h2>
    <div v-if="loadingContracts">
      <AnimationsLoading />
    </div>
    <ul class="mt-2 list-disc pl-12" v-else>
      <li v-for="contract in contracts" class="mt-2">
        <a :href="contract.openseaUrl" target="_blank" class="underline"
          >{{ contract.name }} - {{ contract.blockchain }}
        </a>
      </li>
    </ul>
    <h2 class="text-md ml-2 mt-5 text-gray-700">
      Click sign and a verification message will pop-up in your wallet. We will
      use this signed message to verify your address.
    </h2>
    <div class="ml-2 mb-7 rounded-md">
      <WalletConnected
        @walletChanged="onWalletChanged"
        @walletNotFound="onWalletNotFound"
        :loadingTokenIds="loadingTokenIds"
      />
    </div>
    <div v-if="loadingTokenIds && !searchingWallet">
      <AnimationsLoading />
    </div>
    <div v-if="multipleTokensAlert">
      <hr class="my-8 h-px border-0 bg-gray-200 dark:bg-gray-700" />
      <SectionsMultipleNfts
        :userCollectionsData="tokenIds"
        @assetClicked="onAssetClicked"
      />
      <hr class="my-8 h-px border-0 bg-gray-200 dark:bg-gray-700" />
    </div>
    <div class="mt-8">
      <BadgesError
        v-if="showError"
        :link="errorLink"
        @closeBadgeButtonClicked="() => (showError = false)"
      >
        {{ errorMsg }}</BadgesError
      >
      <BadgesWarning
        v-if="showWarning"
        class="ml-2 mr-10"
        @closeBadgeButtonClicked="() => (showWarning = false)"
        >We couldn't find an asset for the connected wallet.
        <a :href="redirectUri" class="font-bold underline"
          >Go back</a
        ></BadgesWarning
      >
    </div>
    <div class="mt-8 mb-12">
      <ButtonsRounded
        class="w-full"
        :class="[
          !walletConnected ||
          (selectedAsset.id == -1 && multipleTokensAlert) ||
          loadingTokenIds
            ? 'cursor-not-allowed opacity-50'
            : 'opacity-100',
        ]"
        @click="onSignClick"
      >
        Sign
      </ButtonsRounded>
    </div>
  </div>
  <div v-else class="my-20 w-full object-center px-8">
    <h2 class="text-center text-lg">Verifying your signature please wait...</h2>
    <AnimationsLoading class="m-15 mx-auto" />
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import Web3 from "web3";
import { useRoute } from "vue-router";
import { faMasksTheater } from "@fortawesome/free-solid-svg-icons";
const config = useRuntimeConfig();
const idpBackUrl = config.public.IDP_SERVER_URL;
const metamaskUrl = "https://metamask.io/download/";
// var nonce = getCookie("nonce");
// let nonce = useRequestHeaders(["nonce"]);
const route = useRoute();
const nonce = route.query["nonce"];
const redirectUri = route.query["redirect_uri"];
const chainIdMissmatch = -32603;
const chainIdNotFound = 4902;
const loadingTokenIds = ref(true);
const loadingContracts = ref(true);
const isVerifyingSignature = ref(false);
const walletConnected = ref(false);
const searchingWallet = ref(true);
const showWarning = ref(false);
console.log(nonce);
if (process.client) {
  window.sessionStorage.setItem("nonce", nonce);
}

var accounts;
var contracts;
var tokenIds;
var cancelUrl;

//state
const multipleTokensAlert = ref(false);
const selectedAsset = ref({ id: -1, collection: "", blockchain: "" });
const showError = ref(false);
var errorMsg;
var walletAddress;
var errorLink = "";

getPolicyContracts();

function onWalletNotFound() {
  errorLink = metamaskUrl;
  errorMsg = "Wallet not found. Please install it here: ";
  showError.value = true;
}

function onWalletChanged(address, chainId) {
  if (process.client) {
    walletAddress = address;
    walletConnected.value = true;
    searchingWallet.value = false;
    selectedAsset.value.id = -1;
    msgParams.domain.chainId = BigInt(chainId).toString();
    console.log("Wallet changed: " + chainId);
    getTokenIds(address);
  }
}

function onAssetClicked(id, collection, blockchain) {
  selectedAsset.value.id = id;
  selectedAsset.value.collection = collection;
  selectedAsset.value.blockchain = blockchain;
  let collectionAddress = getScAddress(collection, blockchain);
  msgParams.message["asset ID"] = id;
  msgParams.message["asset blockchain"] = blockchain;
  msgParams.message["asset collection"] = collectionAddress;
}

function getScAddress(collection, blockchain) {
  for (const assets of tokenIds) {
    console.log("Name - blockchain: " + assets.name + "-" + assets.blockchain);
    console.log("Targets: " + collection + "-" + blockchain);
    if (
      assets.name === collection &&
      assets.blockchain === blockchain.toLowerCase()
    ) {
      return assets.address;
    }
  }
  return null;
}

function getPolicyContracts() {
  fetch(idpBackUrl + "/api/oidc/getPolicyContracts", {
    method: "GET",
    mode: "cors",
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      } else {
        console.error("Error " + response.status);
      }
    })
    .then((data) => {
      contracts = data.map((contractInfo) => {
        const openseaUrl =
          "https://testnets.opensea.io/assets/goerli/" + contractInfo.address;
        const name = contractInfo.name;
        const blockchain =
          contractInfo.blockchain.charAt(0).toUpperCase() +
          contractInfo.blockchain.slice(1);
        return { openseaUrl, name, blockchain };
      });
      console.log("Contracts: " + JSON.stringify(contracts));
    })
    .finally(() => {
      loadingContracts.value = false;
    });
}

function getTokenIds(walletAddress) {
  showWarning.value = false;
  loadingTokenIds.value = true;
  multipleTokensAlert.value = false;
  fetch(idpBackUrl + "/api/oidc/assets" + "?address=" + walletAddress, {
    method: "GET",
    mode: "cors",
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      } else {
        console.error(error);
      }
    })
    .then((data) => {
      if (data == null) {
        cancelUrl = idpBackUrl + "/api/oidc/cancel" + "?nonce=" + nonce;
        showWarning.value = true;
      } else {
        console.log("DATA: " + JSON.stringify(data));
        tokenIds = data;
        multipleTokensAlert.value = true;
      }
    })
    .finally(() => {
      loadingTokenIds.value = false;
    });
}

function onSignClick() {
  if (
    (selectedAsset.value.id == -1 && multipleTokensAlert.value) ||
    !walletConnected.value ||
    loadingTokenIds.value
  ) {
    if (selectedAsset.value.id == -1 && multipleTokensAlert.value) {
      errorMsg = "Please select an asset before signing the message";
      showError.value = true;
    } else if (!walletConnected.value) {
      console.log("Wallet not connected");
      errorMsg = "Please connect your wallet to sign the message";
      showError.value = true;
    }
  } else {
    console.log("Signature started");
    var params = [walletAddress, JSON.stringify(msgParams)];
    var method = "eth_signTypedData_v4";
    window.ethereum
      .request({
        method: method,
        params: params,
      })
      .then(onSigned, onError);
  }
}

async function onError(error) {
  console.log(error);
  if (error.code == chainIdMissmatch) {
    console.log("MISMATCH");
    try {
      await swicthEthereumChain();
      onSignClick(); //If the chain is switched successfully we send the sign operation again
    } catch (error) {
      console.log("User refused to sign/switch network");
    }
  }
}

async function onSigned(result) {
  let signatureJSON = {
    message: msgParams,
    signature: result,
  };
  console.log(`SIGNATURE JSON: ${JSON.stringify(signatureJSON)}`);
  isVerifyingSignature.value = true;
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
        return response.json();
      } else {
        console.log("Error: ", response.status);
        if (response.status >= 500) {
          return "Something went wrong. Please try again later";
        }
        return response.text();
      }
    })
    .then((parsedResponse) => {
      if (typeof parsedResponse === "string") {
        console.log("Error description: " + parsedResponse);
        window.location.href = redirectUri + "?error=" + parsedResponse;
      } else {
        window.location.href =
          redirectUri +
          "?code=" +
          parsedResponse.code +
          "&state=" +
          parsedResponse.state +
          "&nonce=" +
          nonce;
      }
    });
}

async function swicthEthereumChain() {
  try {
    await window.ethereum.request({
      method: "wallet_switchEthereumChain",
      params: [{ chainId: config.public.CHAIN_ID }],
    });
  } catch (error) {
    console.log("ERROR switching the user's network: " + error);
    if (error.code === chainIdNotFound) {
      try {
        await window.ethereum.request({
          method: "wallet_addEthereumChain",
          params: [
            {
              chainId: config.public.CHAIN_ID,
              rpcUrl: config.public.RPC_URL,
            },
          ],
        });
        await switchEthereumChain();
      } catch (addError) {
        console.error(addError);
        errorMsg = "Failed adding the required network, please add it manually";
        errorMsg.value = true;
        throw new Error("Network not found");
      }
    } else {
      errorMsg = "Signature rejected";
      showError.value = true;
      throw new Error("User refused to switch networks");
    }
    console.error(error);
  }
}

function getCookie(name) {
  // Split cookie string and get all individual name=value pairs in an array
  if (typeof window !== "undefined") {
    console.log(document.cookie);
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
    chainId: BigInt(5).toString(),
    name: "NFT ownership verification service",
    verifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
    version: "1",
  },
  message: {
    // why: "To prove you have a NFT of one of the mentioned above collections, please sign this message",
    timestamp: getTimeAndDate(),
    nonce: nonce,
    "asset ID": 0,
    "asset collection": "",
  },
  // Refers to the keys of the *types* object below.
  primaryType: "Entrust Ownership Based Access Control as a Service",
  types: {
    // Clarify if EIP712Domain refers to the domain the contract is hosted on
    EIP712Domain: [
      { name: "name", type: "string" },
      { name: "version", type: "string" },
      { name: "chainId", type: "uint256" },
      { name: "verifyingContract", type: "address" },
    ],

    "Entrust Ownership Based Access Control as a Service": [
      { name: "timestamp", type: "string" },
      { name: "nonce", type: "string" },
      { name: "asset ID", type: "uint256" },
      { name: "asset collection", type: "address" },
      { name: "asset blockchain", type: "string" },
    ],
  },
};
</script>
