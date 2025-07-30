// import { createWeb3Modal, defaultConfig } from '@web3modal/ethers/vue'
// import { Contract, BrowserProvider } from 'ethers'
// import { bsc } from 'viem/chains'
// import { $modals } from '../modals'

// const projectId = 'a6a3ec5f1992123e14cbed3969bc202a'

// const chain = {
//   chainId: bsc.id,
//   name: bsc.name,
//   currency: bsc.nativeCurrency.symbol,
//   explorerUrl: bsc.blockExplorers.default.url,
//   rpcUrl: bsc.rpcUrls.default.http[0],
// }

// const metadata = {
//   name: 'Web3Modal',
//   description: 'Web3Modal Example',
//   url: location.origin,
//   icons: ['https://avatars.githubusercontent.com/u/37784886'],
// }

// let modal = null

// function getConnectProvider() {
//   if (!modal) {
//     modal = createWeb3Modal({
//       ethersConfig: defaultConfig({ metadata }),
//       chains: [chain],
//       projectId,
//       enableAnalytics: true,
//     })
//   }

//   return modal
// }

// const walletInterface = {
//   open() {
//     getConnectProvider().open()
//   },
//   close() {
//     getConnectProvider().close()
//   },
//   disconect() {
//     getConnectProvider().disconnect()
//   },
//   getWalletProvider() {
//     return getConnectProvider().getWalletProvider()
//   },
//   async isConnected() {
//     const provider = this.getWalletProvider()

//     if (!provider) {
//       walletInterface.open()
//       return null
//     }

//     if (getConnectProvider().getChainId() !== chain.chainId) {
//       await getConnectProvider().switchNetwork(chain.chainId)
//     }

//     return provider
//   },
// }

// const etherInterface = {
//   async ethersProvider() {
//     if (await walletInterface.isConnected()) {
//       return new BrowserProvider(walletInterface.getWalletProvider())
//     }
//     return null
//   },
//   async etherGetSigner() {
//     if (await walletInterface.isConnected()) {
//       return (await this.ethersProvider()).getSigner()
//     }
//     return null
//   },
// }

// async function createContract(address: string, abi: any[]) {
//   if (!address || !abi) return
//   if (!(await $wallet.isConnected())) return Promise.reject({ message: 'Wallet not connected' })
//   try {
//     const signer = await $wallet.etherGetSigner()
//     return new Contract(address, abi, signer)
//   } catch (error) {
//     $modals.error.show({
//       text: JSON.stringify(error),
//     })
//     return Promise.reject(error)
//   }
// }

// // async function getPayMethod(): Promise<DistRes | BuyRes> {
// //   if (data.code === 'distributor') {
// //     return await $requests.transactions.setBuyDistributor()
// //   }

// //   return await $requests.transactions.setBuy({ code: data.code })
// // }

// // async function pay() {
// //   if (loading.value) return

// //   loading.value = true
// //   status.value = 'Loading'

// //   try {
// //     await $store.updateCfg()

// //     const balance = await $contracts.usdtBnb.getBalanceOf()
// //     const amount = BigInt($formatAmountToEther(data.amount))

// //     if (balance <= BigInt(0) || balance < amount) {
// //       showError($t('PaymentsMessages.InsufficientFunds'))

// //       loading.value = false
// //       return
// //     }

// //     if (balance >= amount) {
// //       const allowance = await $contracts.usdtBnb.getAllowance()

// //       if (allowance < amount) {
// //         const approve = await $contracts.usdtBnb.setApprove(amount - allowance)
// //         await approve.wait()
// //       }

// //       const { uid } = await getPayMethod()

// //       if (!uid) {
// //         loading.value = false
// //         showError($t('PaymentsMessages.failed'))
// //         return
// //       }

// //       const deposit = await $contracts.main.setDeposit(amount, uid)
// //       await deposit.wait()

// //       showSuccess($t('PaymentsMessages.success'))

// //       $modals.paySystem.close()
// //     }
// //   } catch (error) {
// //     if (!error) return
// //     if ($isApiError(error)) {
// //       if (error.message) {
// //         showError(JSON.stringify(error.message))
// //       } else {
// //         showError(JSON.stringify(error))
// //       }
// //     }
// //   }

// //   loading.value = false
// // }

// // function generatePCQR() {
// //   const options = {
// //     quality: 0.3,
// //     margin: 1,
// //     width: 300,
// //     color: {
// //       dark: $lchToHex($themeColors.baseContent),
// //       light: $lchToHex($themeColors.base100),
// //     },
// //   }
// //   QRCode.toCanvas(qrzone.value, getText(), options)
// // }

// const $wallet = {
//   ...walletInterface,
//   ...etherInterface,
//   createContract,
// }

// export { $wallet }
