import { useState, useMemo, useCallback, useEffect } from 'react';
import {
  TextField,
  Box,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Typography,
  CircularProgress,
  InputAdornment,
  IconButton,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import ClearIcon from '@mui/icons-material/Clear';
import debounce from 'lodash.debounce';
import { userApi } from '../../api';
import type { User } from '../../api/types';

interface UserSearchProps {
  selectedUser: User | null;
  onSelect: (user: User | null) => void;
}

const UserSearch = ({ selectedUser, onSelect }: UserSearchProps) => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [open, setOpen] = useState(false);

  const debouncedSearch = useMemo(
    () =>
      debounce(async (searchQuery: string) => {
        if (!searchQuery.trim()) {
          setResults([]);
          setLoading(false);
          return;
        }
        try {
          const { data } = await userApi.search(searchQuery);
          setResults(data);
        } catch {
          setResults([]);
        }
        setLoading(false);
      }, 300),
    []
  );

  useEffect(() => {
    return () => {
      debouncedSearch.cancel();
    };
  }, [debouncedSearch]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setQuery(value);
    setLoading(true);
    setOpen(true);
    debouncedSearch(value);
  };

  const handleSelect = (user: User) => {
    onSelect(user);
    setQuery('');
    setResults([]);
    setOpen(false);
  };

  const handleClear = () => {
    onSelect(null);
    setQuery('');
    setResults([]);
    setOpen(false);
  };

  return (
    <Box sx={{ position: 'relative' }}>
      <TextField
        fullWidth
        label="Assignee"
        placeholder={selectedUser ? '' : 'Search by name or email'}
        value={selectedUser ? selectedUser.name : query}
        onChange={handleInputChange}
        disabled={!!selectedUser}
        onFocus={() => setOpen(true)}
        slotProps={{
          input: {
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon />
              </InputAdornment>
            ),
            endAdornment: selectedUser ? (
              <InputAdornment position="end">
                <IconButton size="small" onClick={handleClear}>
                  <ClearIcon />
                </IconButton>
              </InputAdornment>
            ) : loading ? (
              <InputAdornment position="end">
                <CircularProgress size={20} />
              </InputAdornment>
            ) : null,
          },
        }}
      />
      
      {open && (results.length > 0 || (!loading && query && !selectedUser)) && (
        <Box
          sx={{
            position: 'absolute',
            top: '100%',
            left: 0,
            right: 0,
            zIndex: 1000,
            bgcolor: 'background.paper',
            border: '1px solid',
            borderColor: 'divider',
            borderRadius: 1,
            maxHeight: 200,
            overflow: 'auto',
            boxShadow: 2,
          }}
        >
          {loading ? (
            <Box sx={{ p: 2, textAlign: 'center' }}>
              <CircularProgress size={20} />
            </Box>
          ) : results.length > 0 ? (
            <List dense>
              {results.map((user) => (
                <ListItem key={user.id} disablePadding>
                  <ListItemButton onClick={() => handleSelect(user)}>
                    <ListItemText
                      primary={user.name}
                      secondary={user.email}
                    />
                  </ListItemButton>
                </ListItem>
              ))}
            </List>
          ) : query && !selectedUser ? (
            <Box sx={{ p: 2, textAlign: 'center' }}>
              <Typography variant="body2" color="text.secondary">
                No users found
              </Typography>
            </Box>
          ) : null}
        </Box>
      )}
    </Box>
  );
};

export default UserSearch;
