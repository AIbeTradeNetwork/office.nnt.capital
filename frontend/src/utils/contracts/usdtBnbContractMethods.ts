// import mainContractData from './mainContractData'
// import USDTContractData from './usdtBnbContractData'
// import { $modals } from 'utils/modals'
// import { $wallet } from 'utils/wallets/web3connect'

// async function getAllowance() {
//   if (!(await $wallet.isConnected())) return Promise.reject({ message: 'Wallet not connected' })

//   try {
//     const signer = (await $wallet.etherGetSigner()).address
//     const contract = await $wallet.createContract(USDTContractData.address, USDTContractData.abi)
//     const data = await contract.allowance(signer, mainContractData.address)

//     return data
//   } catch (error) {
//     $modals.error.show({
//       text: JSON.stringify(error),
//     })
//     return Promise.reject(error)
//   }
// }

// async function getBalanceOf() {
//   if (!(await $wallet.isConnected())) return Promise.reject({ message: 'Wallet not connected' })
//   try {
//     const signer = (await $wallet.etherGetSigner()).address
//     const contract = await $wallet.createContract(USDTContractData.address, USDTContractData.abi)
//     const data = await contract.balanceOf(signer)
//     return data
//   } catch (error) {
//     $modals.error.show({
//       text: JSON.stringify(error),
//     })
//     return Promise.reject(error)
//   }
// }

// async function setApprove(amount: bigint) {
//   if (!(await $wallet.isConnected())) return Promise.reject({ message: 'Wallet not connected' })
//   try {
//     const contract = await $wallet.createContract(USDTContractData.address, USDTContractData.abi)
//     const data = await contract.approve(mainContractData.address, amount)
//     return data
//   } catch (error) {
//     $modals.error.show({
//       text: JSON.stringify(error),
//     })
//     return Promise.reject(error)
//   }
// }

// const usdtBnbContractMethods = {
//   getAllowance,
//   getBalanceOf,
//   setApprove,
// }

// export { usdtBnbContractMethods }
