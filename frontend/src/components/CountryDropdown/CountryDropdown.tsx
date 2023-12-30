import React, { FC, KeyboardEvent } from 'react';
import './CountryDropdown.scss';
import TextField from '@mui/material/TextField';
import Autocomplete from '@mui/material/Autocomplete';

interface CountryDropdownProps {
  countries: string[],
  send: (event: KeyboardEvent<HTMLInputElement>, role: string) => void,
  isDisabled: boolean
}

const CountryDropdown: FC<CountryDropdownProps> = ({countries, send, isDisabled}) => (
    <Autocomplete
      disablePortal
      id="country-dropdown"
      className="CountryDropdown"
      disabled={isDisabled}
      options={countries}
      renderInput={(params) => <TextField {...params} label="Select a destination" />}
      onKeyDown={(e) => (e.key === "Enter") ? send(e as React.KeyboardEvent<HTMLInputElement>, "unknown") : undefined}
    />
);

export default CountryDropdown;