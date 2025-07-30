// import contractData from './mainContractData'
// import { $modals } from 'utils/modals'
// import { $wallet } from 'utils/wallets/web3connect'
// // const { encodeBytes32String } = await import('ethers')

// async function getUsdtTokenTest() {
//   if (!(await $wallet.isConnected())) return Promise.reject({ message: 'Wallet not connected' })

//   try {
//     const contract = await $wallet.createContract(contractData.address, contractData.abi)
//     const data = await contract.usdtToken()
//     return data
//   } catch (error) {
//     $modals.error.show({
//       text: JSON.stringify(error),
//     })
//     return Promise.reject(error)
//   }
// }

// async function setDeposit(amount: bigint, uuid: string) {
//   if (!(await $wallet.isConnected())) return Promise.reject({ message: 'Wallet not connected' })

//   try {
//     const contract = await $wallet.createContract(contractData.address, contractData.abi)
//     // const uidBite32 = encodeBytes32String(uuid)
//     const data = await contract.deposit(amount, uuid)
//     return data
//   } catch (error) {
//     $modals.error.show({
//       text: JSON.stringify(error),
//     })
//     return Promise.reject(error)
//   }
// }

// const mainContractMethods = {
//   getUsdtTokenTest,
//   setDeposit,
// }

// export { mainContractMethods }
