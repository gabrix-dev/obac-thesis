<template>
  <div v-if="searchingWallet">
      <AnimationsLoading />
    </div>
  <div v-else-if="walletConnected">
    <div  class="mt-8 flex justify-between shadow-md px-6 py-4 rounded-lg items-center">
        <p class="flex items-center">
          <IconsTick />
          <h2 class="ml-3 text-lg font-bold text-gray-800">Wallet connected:</h2>
          <span class="ml-5 mr-10 text-gray-700">{{ wallet }}</span>
        </p>
        <button
          class="items-center justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 hover:bg-indigo-400"
          @click="connectWallet"
        >
          Change Wallet
        </button>
      </div>
  </div>
  <div v-else>
      <div class="mt-8 flex justify-between shadow-md px-6 py-4 rounded-lg items-center">
        <p class="flex items-center">
          <IconsCross />
          <h2 class="ml-3 text-lg font-bold text-gray-800">
            Wallet not connected
          </h2>
        </p>
        <button
          class="items-center justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 hover:bg-indigo-400"
          @click="connectWallet"
        >
          Connect Wallet
        </button>
      </div>
    </div>
</template>

<script setup>
const emit = defineEmits(["walletChanged", "walletNotFound"]);
const searchingWallet = ref(true);
const walletConnected = ref(false);
const wallet = ref("");
var walletAddress;

checkIfWalletConnected();

if (process.client && window.ethereum ){
  window.ethereum.on('chainChanged', handleChainChanged);
}

function handleChainChanged(id) {
  console.log("Chain ID changed: "+id)
  emit("walletChanged", walletAddress, id);
}

async function checkIfWalletConnected() {
  if (process.client) {
    console.log("Checking if wallet connected")
    if (window.ethereum) {
      console.log("Checking if wallet connected")
      try {
        const response = await window.ethereum.request({
          method: "eth_requestAccounts",
        });
        walletAddress =response[0]
        wallet.value = getShortnedAddressStr(walletAddress);
        walletConnected.value = true;
        const chainId = await window.ethereum.request({
          method: "eth_chainId",
        });
        console.log("Wallet changed, chainid - address: "+chainId+" - "+walletAddress)
        emit("walletChanged", walletAddress, chainId);
      } catch (error) {
        console.log(error);
      }
    }
    searchingWallet.value = false;
  }
}

async function connectWallet() {
  if (process.client) {
    if (window.ethereum) {
      try {
        await ethereum.request({
          method: "wallet_requestPermissions",
          params: [{ eth_accounts: {} }],
        });

        const response = await window.ethereum.request({
          method: "eth_requestAccounts",
        });
        const chainId = await window.ethereum.request({
          method: "eth_chainId",
        });
        walletAddress = response[0]
        wallet.value = getShortnedAddressStr(walletAddress);
        console.log(wallet.value);
        walletConnected.value = true;
        emit("walletChanged", walletAddress, chainId);
      } catch (error) {
        console.log(error);
      }
    } else {
      console.log("Wallet not found");
      emit("walletNotFound");
    }
  }
}

function getShortnedAddressStr(hexString){
  return hexString.slice(0, 6) + "......" + hexString.slice(-4)
}
</script>
