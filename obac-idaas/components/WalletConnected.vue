<template>
  <div>
    <div v-if="searchingWallet">
      <AnimationsLoading />
    </div>
    <div v-else-if="walletConnected">
      <div class="mt-8 flex items-center justify-between mr-10 border rounded-lg p-3 border-gray-300">
        <div class="flex items-center">
          <IconsPurpleTick class="shrink-0"/>
          <h2 class="ml-3 text-lg font-bold text-gray-800">Wallet connected:</h2>
          <span class="ml-5 mr-10 text-gray-700 shrink">{{ wallet }}</span>
</div>
        <button
          class="ml-8 rounded-md bg-entrust-purple px-4 py-2 pl-2 text-sm font-medium text-white transition-colors focus:outline-none focus:ring-2 focus:ring-purple-400 focus:ring-offset-2 hover:bg-entrust-dark-purple"
          :class="[
          loadingTokenIds
            ? 'cursor-not-allowed opacity-50'
            : 'opacity-100',
        ]"
          @click="connectWallet"
        >
          Change Wallet
        </button>
      </div>
    </div>
    <div v-else>
      <div class="mt-8 flex justify-between items-center mr-20">
        <p class="flex items-center">
          <IconsCross />
          <h2 class="ml-3 text-lg font-bold text-gray-800">
            Wallet not connected
          </h2>
        </p>
        <button
          class="ml-8 rounded-md bg-entrust-purple px-4 py-2 pl-2 text-sm font-medium text-white transition-colors focus:outline-none focus:ring-2 focus:ring-purple-400 focus:ring-offset-2 hover:bg-entrust-dark-purple"
          @click="connectWallet"
        >
          Connect Wallet
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
const emit = defineEmits(["walletChanged", "walletNotFound"]);
const searchingWallet = ref(true);
const walletConnected = ref(false);
const wallet = ref("");
const props = defineProps({
  loadingTokenIds: {
    type: Boolean,
    required: false,
  }
});

checkIfWalletConnected();


if (process.client && window.ethereum){
  window.ethereum.on('chainChanged', handleChainChanged);
}

function handleChainChanged(id) {
  console.log("Chain ID changed: "+id)
  emit("walletChanged", wallet.value, id);
}

async function checkIfWalletConnected2() {
  if (process.client) {
      console.log("Checking if wallet connected")

    if (window.ethereum) {
      console.log("Checking if wallet connected")
      await window.ethereum.request({
        method: "eth_requestAccounts",
      })
      .then((response) => { 
        console.log("Response: ",response)
        wallet.value = getShortnedAddressStr(response[0]);
        emit("walletChanged", response[0]);
        walletConnected.value = true;
      })
      .catch((error) => console.log(error));
    }
    searchingWallet.value = false;
  }
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
        let walletAddress =response[0]
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
  console.log("LOADING: "+props.loadingTokenIds)
  if (process.client) {
    if (window.ethereum ) {
      const chainId = await window.ethereum.request({
          method: "eth_chainId",
        });
      if(!props.loadingTokenIds){
      await ethereum.request({
        method: "wallet_requestPermissions",
        params: [{ eth_accounts: {} }],
      });
      await window.ethereum.request({
        method: "eth_requestAccounts",
      })
      .then((response) => { 
        wallet.value = getShortnedAddressStr(response[0]);
        console.log(wallet.value);
        walletConnected.value = true;
        emit("walletChanged", response[0], chainId);
      })
      .catch((error) => console.log(error));
    }
     } else {
      console.log("Wallet not found");
      emit("walletNotFound")
      }
    }
  }

  function getShortnedAddressStr(hexString){
    return hexString.slice(0, 6) + "......" + hexString.slice(-4)
}

</script>
