
pragma solidity >=0.7.0 <0.9.0;


contract Hello {

    uint256 number;
    string public message;

    /**
     * @dev Store value in variable
     * @param num value to store
     */
    function store(uint256 num) public {
        number = num;
    }

    /**
     * @dev Return value 
     * @return value of 'number'
     */
    function retrieve() public view returns (uint256){
        return number;
    }
    
    
    
     function setMessage(string memory newMessage) public
     {
         message=newMessage;
     }

    function getMessage()public view returns(string memory)
    {
        return message;
    }
    
}
