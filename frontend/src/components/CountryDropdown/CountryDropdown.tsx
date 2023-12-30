import React, { FC, KeyboardEvent } from 'react';
import './CountryDropdown.scss';
import TextField from '@mui/material/TextField';
import Autocomplete from '@mui/material/Autocomplete';
import { MESSAGE_TYPE_COUNTRY, VALID_MESSAGE_TYPE } from '../../api/websocket';

interface CountryDropdownProps {
  countries: string[],
  send: (event: KeyboardEvent<HTMLInputElement>, role: string, messageType: VALID_MESSAGE_TYPE) => void,
  isDisabled: boolean
}

const CountryDropdown: FC<CountryDropdownProps> = ({countries, send, isDisabled}) => {
    return <Autocomplete
      disablePortal
      id="country-dropdown"
      className="CountryDropdown"
      disabled={isDisabled}
      options={countries}
      renderInput={(params) => <TextField {...params} label="Select a destination" />}
      onKeyDown={(e) => (e.key === "Enter") ? send(e as React.KeyboardEvent<HTMLInputElement>, "SET_COUNTRY", MESSAGE_TYPE_COUNTRY) : undefined}
    />
};

export default CountryDropdown;