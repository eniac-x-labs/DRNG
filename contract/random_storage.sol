// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

// solc --abi ./contract/random_storage.sol -o ./contract --overwrite
// abigen --abi contract/RandomHashListStorageContract.abi --pkg contract --type RandomHashListStorage --out contract/random_hash_list_storage.go

contract RandomHashListStorageContract {
    mapping(string => mapping(address => string[])) private reqIDToAddressToCollection;
    mapping(string => string[]) private hashLists;
    mapping(string => address[]) private reqIDToAddresses;

    event HashListStored(string indexed reqID, string[] hashList);
    event CollectionSaved(string indexed reqID, address indexed sender, string[] collection);

    function storeHashList(string memory reqID, string[] memory _hashList) public {
        hashLists[reqID] = _hashList;
        emit HashListStored(reqID, _hashList);
    }

    function getHashList(string memory reqID) public view returns (string[] memory) {
        return hashLists[reqID];
    }

    function hashExistsInList(string memory reqID, string memory hash) public view returns (bool) {
        string[] memory _hashList = hashLists[reqID];

        for (uint256 i = 0; i < _hashList.length; i++) {
            if (keccak256(abi.encodePacked(_hashList[i])) == keccak256(abi.encodePacked(hash))) {
                return true;
            }
        }
        return false;
    }

    function generatorSaveCollection(string memory reqID, string[] memory collection) public {
        reqIDToAddressToCollection[reqID][msg.sender] = collection;
        reqIDToAddresses[reqID].push(msg.sender);
        emit CollectionSaved(reqID, msg.sender, collection);
    }

    function getCollectionByReqIDAndAddress(string memory reqID, address addr) public view returns (string[] memory) {
        return reqIDToAddressToCollection[reqID][addr];
    }

    function getAddressesByReqID(string memory reqID) public view returns (address[] memory) {
        return reqIDToAddresses[reqID];
    }
}