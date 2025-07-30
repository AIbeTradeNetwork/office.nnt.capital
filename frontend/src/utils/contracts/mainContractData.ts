// export default {
//   address: <`0x${string}`>'0xd14fa0529d587e6806E7c2eC9f57bDA1Cf6d2854',
//   abi: [
//     {
//       inputs: [],
//       stateMutability: 'nonpayable',
//       type: 'constructor',
//     },
//     {
//       inputs: [],
//       name: 'EnforcedPause',
//       type: 'error',
//     },
//     {
//       inputs: [],
//       name: 'ExpectedPause',
//       type: 'error',
//     },
//     {
//       anonymous: false,
//       inputs: [
//         {
//           indexed: false,
//           internalType: 'address',
//           name: 'previousAdmin',
//           type: 'address',
//         },
//         {
//           indexed: false,
//           internalType: 'address',
//           name: 'newAdmin',
//           type: 'address',
//         },
//       ],
//       name: 'AdminChanged',
//       type: 'event',
//     },
//     {
//       anonymous: false,
//       inputs: [
//         {
//           indexed: true,
//           internalType: 'address',
//           name: 'beacon',
//           type: 'address',
//         },
//       ],
//       name: 'BeaconUpgraded',
//       type: 'event',
//     },
//     {
//       anonymous: false,
//       inputs: [
//         {
//           indexed: true,
//           internalType: 'address',
//           name: 'sender',
//           type: 'address',
//         },
//         {
//           indexed: false,
//           internalType: 'uint256',
//           name: 'amount',
//           type: 'uint256',
//         },
//         {
//           indexed: true,
//           internalType: 'string',
//           name: 'uuid',
//           type: 'string',
//         },
//       ],
//       name: 'Deposit',
//       type: 'event',
//     },
//     {
//       anonymous: false,
//       inputs: [
//         {
//           indexed: false,
//           internalType: 'uint8',
//           name: 'version',
//           type: 'uint8',
//         },
//       ],
//       name: 'Initialized',
//       type: 'event',
//     },
//     {
//       anonymous: false,
//       inputs: [
//         {
//           indexed: false,
//           internalType: 'address',
//           name: 'account',
//           type: 'address',
//         },
//       ],
//       name: 'Paused',
//       type: 'event',
//     },
//     {
//       anonymous: false,
//       inputs: [
//         {
//           indexed: true,
//           internalType: 'bytes32',
//           name: 'role',
//           type: 'bytes32',
//         },
//         {
//           indexed: true,
//           internalType: 'bytes32',
//           name: 'previousAdminRole',
//           type: 'bytes32',
//         },
//         {
//           indexed: true,
//           internalType: 'bytes32',
//           name: 'newAdminRole',
//           type: 'bytes32',
//         },
//       ],
//       name: 'RoleAdminChanged',
//       type: 'event',
//     },
//     {
//       anonymous: false,
//       inputs: [
//         {
//           indexed: true,
//           internalType: 'bytes32',
//           name: 'role',
//           type: 'bytes32',
//         },
//         {
//           indexed: true,
//           internalType: 'address',
//           name: 'account',
//           type: 'address',
//         },
//         {
//           indexed: true,
//           internalType: 'address',
//           name: 'sender',
//           type: 'address',
//         },
//       ],
//       name: 'RoleGranted',
//       type: 'event',
//     },
//     {
//       anonymous: false,
//       inputs: [
//         {
//           indexed: true,
//           internalType: 'bytes32',
//           name: 'role',
//           type: 'bytes32',
//         },
//         {
//           indexed: true,
//           internalType: 'address',
//           name: 'account',
//           type: 'address',
//         },
//         {
//           indexed: true,
//           internalType: 'address',
//           name: 'sender',
//           type: 'address',
//         },
//       ],
//       name: 'RoleRevoked',
//       type: 'event',
//     },
//     {
//       anonymous: false,
//       inputs: [
//         {
//           indexed: false,
//           internalType: 'address',
//           name: 'account',
//           type: 'address',
//         },
//       ],
//       name: 'Unpaused',
//       type: 'event',
//     },
//     {
//       anonymous: false,
//       inputs: [
//         {
//           indexed: true,
//           internalType: 'address',
//           name: 'implementation',
//           type: 'address',
//         },
//       ],
//       name: 'Upgraded',
//       type: 'event',
//     },
//     {
//       anonymous: false,
//       inputs: [
//         {
//           indexed: true,
//           internalType: 'address',
//           name: 'recipient',
//           type: 'address',
//         },
//         {
//           indexed: false,
//           internalType: 'uint256',
//           name: 'amount',
//           type: 'uint256',
//         },
//       ],
//       name: 'Withdrawal',
//       type: 'event',
//     },
//     {
//       inputs: [],
//       name: 'DEFAULT_ADMIN_ROLE',
//       outputs: [
//         {
//           internalType: 'bytes32',
//           name: '',
//           type: 'bytes32',
//         },
//       ],
//       stateMutability: 'view',
//       type: 'function',
//     },
//     {
//       inputs: [],
//       name: 'MANAGER_ROLE',
//       outputs: [
//         {
//           internalType: 'bytes32',
//           name: '',
//           type: 'bytes32',
//         },
//       ],
//       stateMutability: 'view',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'uint256',
//           name: '_amount',
//           type: 'uint256',
//         },
//         {
//           internalType: 'string',
//           name: '_uuid',
//           type: 'string',
//         },
//       ],
//       name: 'deposit',
//       outputs: [],
//       stateMutability: 'nonpayable',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'string[]',
//           name: 'uuids',
//           type: 'string[]',
//         },
//       ],
//       name: 'getMultiplePaymentInfo',
//       outputs: [
//         {
//           components: [
//             {
//               internalType: 'uint256',
//               name: 'amount',
//               type: 'uint256',
//             },
//             {
//               internalType: 'uint256',
//               name: 'lastDepositTime',
//               type: 'uint256',
//             },
//           ],
//           internalType: 'struct PaymentGateway.PaymentInfo[]',
//           name: '',
//           type: 'tuple[]',
//         },
//       ],
//       stateMutability: 'view',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'string',
//           name: 'uuid',
//           type: 'string',
//         },
//       ],
//       name: 'getPaymentInfo',
//       outputs: [
//         {
//           components: [
//             {
//               internalType: 'uint256',
//               name: 'amount',
//               type: 'uint256',
//             },
//             {
//               internalType: 'uint256',
//               name: 'lastDepositTime',
//               type: 'uint256',
//             },
//           ],
//           internalType: 'struct PaymentGateway.PaymentInfo',
//           name: '',
//           type: 'tuple',
//         },
//       ],
//       stateMutability: 'view',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'bytes32',
//           name: 'role',
//           type: 'bytes32',
//         },
//       ],
//       name: 'getRoleAdmin',
//       outputs: [
//         {
//           internalType: 'bytes32',
//           name: '',
//           type: 'bytes32',
//         },
//       ],
//       stateMutability: 'view',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'bytes32',
//           name: 'role',
//           type: 'bytes32',
//         },
//         {
//           internalType: 'address',
//           name: 'account',
//           type: 'address',
//         },
//       ],
//       name: 'grantRole',
//       outputs: [],
//       stateMutability: 'nonpayable',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'bytes32',
//           name: 'role',
//           type: 'bytes32',
//         },
//         {
//           internalType: 'address',
//           name: 'account',
//           type: 'address',
//         },
//       ],
//       name: 'hasRole',
//       outputs: [
//         {
//           internalType: 'bool',
//           name: '',
//           type: 'bool',
//         },
//       ],
//       stateMutability: 'view',
//       type: 'function',
//     },
//     {
//       inputs: [],
//       name: 'initialize',
//       outputs: [],
//       stateMutability: 'nonpayable',
//       type: 'function',
//     },
//     {
//       inputs: [],
//       name: 'pause',
//       outputs: [],
//       stateMutability: 'nonpayable',
//       type: 'function',
//     },
//     {
//       inputs: [],
//       name: 'paused',
//       outputs: [
//         {
//           internalType: 'bool',
//           name: '',
//           type: 'bool',
//         },
//       ],
//       stateMutability: 'view',
//       type: 'function',
//     },
//     {
//       inputs: [],
//       name: 'proxiableUUID',
//       outputs: [
//         {
//           internalType: 'bytes32',
//           name: '',
//           type: 'bytes32',
//         },
//       ],
//       stateMutability: 'view',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'bytes32',
//           name: 'role',
//           type: 'bytes32',
//         },
//         {
//           internalType: 'address',
//           name: 'account',
//           type: 'address',
//         },
//       ],
//       name: 'renounceRole',
//       outputs: [],
//       stateMutability: 'nonpayable',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'bytes32',
//           name: 'role',
//           type: 'bytes32',
//         },
//         {
//           internalType: 'address',
//           name: 'account',
//           type: 'address',
//         },
//       ],
//       name: 'revokeRole',
//       outputs: [],
//       stateMutability: 'nonpayable',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'address',
//           name: '_newAddress',
//           type: 'address',
//         },
//       ],
//       name: 'setTokenAddr',
//       outputs: [],
//       stateMutability: 'nonpayable',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'bytes4',
//           name: 'interfaceId',
//           type: 'bytes4',
//         },
//       ],
//       name: 'supportsInterface',
//       outputs: [
//         {
//           internalType: 'bool',
//           name: '',
//           type: 'bool',
//         },
//       ],
//       stateMutability: 'view',
//       type: 'function',
//     },
//     {
//       inputs: [],
//       name: 'unpause',
//       outputs: [],
//       stateMutability: 'nonpayable',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'address',
//           name: 'newImplementation',
//           type: 'address',
//         },
//       ],
//       name: 'upgradeTo',
//       outputs: [],
//       stateMutability: 'nonpayable',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'address',
//           name: 'newImplementation',
//           type: 'address',
//         },
//         {
//           internalType: 'bytes',
//           name: 'data',
//           type: 'bytes',
//         },
//       ],
//       name: 'upgradeToAndCall',
//       outputs: [],
//       stateMutability: 'payable',
//       type: 'function',
//     },
//     {
//       inputs: [],
//       name: 'usdtToken',
//       outputs: [
//         {
//           internalType: 'contract IERC20',
//           name: '',
//           type: 'address',
//         },
//       ],
//       stateMutability: 'view',
//       type: 'function',
//     },
//     {
//       inputs: [
//         {
//           internalType: 'address payable',
//           name: '_addr',
//           type: 'address',
//         },
//         {
//           internalType: 'uint256',
//           name: '_amount',
//           type: 'uint256',
//         },
//       ],
//       name: 'withdraw',
//       outputs: [],
//       stateMutability: 'nonpayable',
//       type: 'function',
//     },
//     {
//       stateMutability: 'payable',
//       type: 'receive',
//     },
//   ],
// }
