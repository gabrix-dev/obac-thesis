<template>
  <Container>
    <!-- <div
      class="mt-12 grid w-full grid-cols-3 justify-between space-x-4"
      :class="'grid-cols-' + tickets.length"
    > -->
    <!-- <section class="g-white dabrk:bg-gray-800"> -->
    <section class="bg-slate-100 bg-opacity-20 backdrop-blur-sm">
      <div class="container mx-auto">
        <div
          class="-mx-6 grid gap-6 sm:grid-cols-2 sm:gap-8 lg:grid-cols-3 xl:grid-cols-3"
        >
          <TicketFull
            @buyButtonClicked="(type) => onBuyButtonClicked(type)"
            v-for="ticket in tickets"
            :ticket="ticket"
            :walletConnected="walletConnected"
          />
        </div>
      </div>
      <WalletConnected
        class="mt-2"
        @walletChanged="onWalletChanged"
        @walletNotFound="onWalletNotFound"
      />
      <BadgesLoading v-if="isProcessingTransaction" class="mt-7"
        >Waiting for buy operation completition. Please wait...</BadgesLoading
      >
      <BadgesSuccess
        v-if="showSuccessMessage"
        class="mt-7"
        @closeBadgeButtonClicked="showSuccessMessage = false"
        >Buy operation completed!
        <a :href="ticketUrl" target="_blank" class="underline"
          >See ticket on Opensea</a
        ></BadgesSuccess
      >
      <!-- <BadgesProcessing
        v-if="isProcessingTransaction"
        class="mt-6"
        @closeBadgeButtonClicked="isProcessingTransaction = false"
        >Buy operation sent to your wallet, click confirm and wait until the
        confirmation message appears</BadgesProcessing
      > -->

      <BadgesError
        v-if="showError"
        class="mt-7"
        :link="errorLink"
        @closeBadgeButtonClicked="() => (showError = false)"
      >
        {{ errorMsg }}</BadgesError
      >
    </section>
  </Container>
</template>

<script setup>
import Web3 from "web3";
import myEpicNftAbi from "@/assets/contracts/MyEpicNFT.json";
var contract, ticketUrl, walletAddress, chainId;
var errorLink = "";
var errorMsg;
const showSuccessMessage = ref(false);
const showError = ref(false);
const walletConnected = ref(false);
const isProcessingTransaction = ref(false);
const chainIdNotFound = 4902;
const userDeniedTx = 4001;
const config = useRuntimeConfig();
const metamaskUrl = "https://metamask.io/download/";
const scMap = {
  [config.public.ETH_CHAIN_ID]: {
    sc: config.public.ETH_SC_ADDRESS,
    opensea: "goerli",
  },
  [config.public.POL_CHAIN_ID]: {
    sc: config.public.POL_SC_ADDRESS,
    opensea: "mumbai",
  },
  [config.public.BSC_CHAIN_ID]: {
    sc: config.public.BSC_SC_ADDRESS,
    opensea: "bsc-testnet",
  },
};

if (process.client && window.ethereum != undefined) {
  // You have a web3 browser! Continue below!
  web3 = new Web3(web3.currentProvider);
} else {
  // Warn the user that they need to get a web3 browser
  // Or install MetaMask, maybe with a nice graphic.
}

function onWalletChanged(address, id) {
  if (process.client) {
    showError.value = false;
    walletConnected.value = true;
    console.log("WALLET CHANGED: " + address + " CHAIN: " + id);
    walletAddress = address;
    chainId = id;
    if (scMap[chainId] != undefined) {
      contract = new web3.eth.Contract(myEpicNftAbi.abi, scMap[chainId].sc);
      console.log("It maps to this SC: " + scMap[chainId].sc);
    }
  }
}

function onWalletNotFound() {
  errorLink = metamaskUrl;
  errorMsg = "Wallet not found. Please install it here: ";
  showError.value = true;
}

async function onBuyButtonClicked(type) {
  if (window.ethereum) {
    showError.value = false;
    if (scMap[chainId] == undefined) {
      errorMsg =
        "Please switch your wallet's network to one of the following: Goerli, Mumbai, Test BSC";
      showError.value = true;
    } else {
      if (!walletAddress) {
        errorMsg = "You have to connect your wallet to buy a ticket";
        showError.value = true;
      } else {
        isProcessingTransaction.value = true;
        switch (type) {
          case "Basic":
            await contract.methods
              .buyBasicTicket()
              .send({ from: walletAddress })
              .then(onTransactionCompleted)
              .catch(onErrorHandler);
            break;
          case "Basic plus":
            await contract.methods
              .buyBasicPlusTicket()
              .send({ from: walletAddress })
              .then(onTransactionCompleted)
              .catch(onErrorHandler);
            break;
          default:
            await contract.methods
              .buyPremiumTicket()
              .send({ from: walletAddress })
              .then(onTransactionCompleted)
              .catch(onErrorHandler);
            break;
        }
      }
    }
  }
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
        throw new Error(addError);
      }
    } else {
      errorMsg = "Please switch the wallet's network";
      isProcessingTransaction.value = false;
      showError.value = true;
      throw new Error(errorMsg);
    }
    console.error(error);
  }
}

function onTransactionCompleted(receipt) {
  isProcessingTransaction.value = false;
  ticketUrl =
    "https://testnets.opensea.io/assets/" +
    scMap[chainId].opensea +
    "/" +
    scMap[chainId].sc +
    "/" +
    receipt.events.Transfer.returnValues.tokenId;
  showSuccessMessage.value = true;
  console.log("TOKEN id: " + receipt.events.Transfer.returnValues.tokenId);
}

function onErrorHandler(error) {
  console.log("Something went wrong: " + error.code);
  if (error.code == userDeniedTx) {
    isProcessingTransaction.value = false;
    errorMsg =
      "Transaction rejected. To buy a ticket please accept the transaction.";
    showError.value = true;
  }
}

const tickets = [
  {
    ticketType: "Basic",
    ticketDescription: "Your standard ticket for Qatar 2022",
    ticketPrice: "200€",
    ticketBuyUrl:
      "https://testnets.opensea.io/assets/goerli/0x14e2ec6b42855e15eaf980b61f9a0a56fa002bbf/1",
    ticketBulletPoints: [
      "Seat: normal",
      "Bookable zones: 3 and 2",
      "Dinner: not included",
    ],
  },
  {
    ticketType: "Basic plus",
    ticketDescription: "Upgrade your basic experience",
    ticketPrice: "500€",
    ticketBuyUrl:
      "https://testnets.opensea.io/assets/goerli/0x14e2ec6b42855e15eaf980b61f9a0a56fa002bbf/2",
    ticketBulletPoints: [
      "Seat: normal",
      "Bookable zones: 1",
      "Dinner: snack and drink",
    ],
  },
  {
    ticketType: "VIP premium",
    ticketDescription: "Live an exclusive experience ",
    ticketPrice: "2000€",
    ticketBuyUrl:
      "https://testnets.opensea.io/assets/goerli/0x14e2ec6b42855e15eaf980b61f9a0a56fa002bbf/3",
    ticketBulletPoints: [
      "Seat: premium",
      "Bookable zones: VIP",
      "Dinner: premium catering",
    ],
  },
];
</script>
