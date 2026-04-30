import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Skeleton,
  Box,
  Button,
  Typography,
  Chip,
} from '@mui/material'
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material'
import type { Config } from '../../../apis/proto/opengate/v1/config'

interface ConfigTableProps {
  configs: Config[]
  loading: boolean
  loadingMore: boolean
  hasMore: boolean
  onRowClick: (config: Config) => void
  onEdit: (config: Config) => void
  onDelete: (config: Config) => void
  onLoadMore: () => void
}

export const ConfigTable = ({
  configs,
  loading,
  loadingMore,
  hasMore,
  onRowClick,
  onEdit,
  onDelete,
  onLoadMore,
}: ConfigTableProps) => {
  if (loading && configs.length === 0) {
    return (
      <TableContainer component={Paper} sx={{ borderRadius: 2 }}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Path Prefix</TableCell>
              <TableCell>Target URL</TableCell>
              <TableCell>Auth</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {Array.from({ length: 5 }).map((_, idx) => (
              <TableRow key={idx}>
                <TableCell><Skeleton /></TableCell>
                <TableCell><Skeleton /></TableCell>
                <TableCell><Skeleton /></TableCell>
                <TableCell><Skeleton width={60} /></TableCell>
                <TableCell><Skeleton width={80} /></TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    )
  }

  if (!loading && configs.length === 0) {
    return (
      <Paper sx={{ p: 4, textAlign: 'center', borderRadius: 2 }}>
        <Typography color="text.secondary">
          No configs found. Create your first config to get started.
        </Typography>
      </Paper>
    )
  }

  return (
    <Box>
      <TableContainer component={Paper} sx={{ borderRadius: 2 }}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Path Prefix</TableCell>
              <TableCell>Target URL</TableCell>
              <TableCell>Auth</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {configs.map((config) => (
              <TableRow
                key={config.id}
                hover
                sx={{ cursor: 'pointer' }}
                onClick={() => onRowClick(config)}
              >
                <TableCell>
                  <Typography fontWeight={500}>{config.name}</Typography>
                </TableCell>
                <TableCell>
                  <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                    {config.pathPrefix}
                  </Typography>
                </TableCell>
                <TableCell>
                  <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                    {config.targetUrl}
                  </Typography>
                </TableCell>
                <TableCell>
                  {config.authentication?.required ? (
                    <Chip label="Required" size="small" color="warning" />
                  ) : (
                    <Chip label="None" size="small" variant="outlined" />
                  )}
                </TableCell>
                <TableCell align="right">
                  <IconButton
                    size="small"
                    onClick={(e) => {
                      e.stopPropagation()
                      onEdit(config)
                    }}
                  >
                    <EditIcon fontSize="small" />
                  </IconButton>
                  <IconButton
                    size="small"
                    color="error"
                    onClick={(e) => {
                      e.stopPropagation()
                      onDelete(config)
                    }}
                  >
                    <DeleteIcon fontSize="small" />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {hasMore && (
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
          <Button
            variant="outlined"
            onClick={onLoadMore}
            disabled={loadingMore}
          >
            {loadingMore ? 'Loading...' : 'Load More'}
          </Button>
        </Box>
      )}
    </Box>
  )
}
