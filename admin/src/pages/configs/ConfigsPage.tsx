import { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Box, Container, Button, TextField, InputAdornment } from '@mui/material'
import { Add as AddIcon, Search as SearchIcon } from '@mui/icons-material'
import { ConfirmDialog, useNotification } from '@gofreego/tsutils'
import { useConfigs } from '../../hooks/useConfigs'
import { ConfigTable } from './components/ConfigTable'
import { ConfigFormDialog } from './components/ConfigFormDialog'
import { ConfigViewDialog } from './components/ConfigViewDialog'
import { PageHeader } from '../../components'
import type { Config, CreateConfigRequest, UpdateConfigRequest } from '../../apis/proto/opengate/v1/config'

export const ConfigsPage = () => {
  const [searchParams, setSearchParams] = useSearchParams()

  const {
    configs,
    loading,
    loadingMore,
    selectedConfig,
    hasMore,
    loadConfigs,
    createConfig,
    updateConfig,
    deleteConfig,
    getConfigById,
    clearSelectedConfig,
  } = useConfigs()

  const { showNotification } = useNotification()

  const [openFormDialog, setOpenFormDialog] = useState(false)
  const [openViewDialog, setOpenViewDialog] = useState(false)
  const [openConfirmDialog, setOpenConfirmDialog] = useState(false)
  const [deleteId, setDeleteId] = useState<string>('')
  const [editData, setEditData] = useState<Config | null>(null)
  const [searchInput, setSearchInput] = useState<string>(searchParams.get('search') || '')
  const [searchTerm, setSearchTerm] = useState<string>(searchParams.get('search') || '')

  const LIMIT = 10

  // Sync filters with URL
  useEffect(() => {
    const params = new URLSearchParams()
    if (searchTerm) params.set('search', searchTerm)
    else params.delete('search')
    setSearchParams(params, { replace: true })
  }, [searchTerm, setSearchParams])

  // Debounce inputs into filters
  useEffect(() => {
    const timer = setTimeout(() => {
      setSearchTerm(searchInput)
    }, 400)
    return () => clearTimeout(timer)
  }, [searchInput])

  // Load initial configs with pagination and filters
  useEffect(() => {
    loadConfigs({ limit: LIMIT, offset: 0, search: searchTerm, append: false })
  }, [searchTerm, loadConfigs])

  const handleLoadMore = () => {
    loadConfigs({ limit: LIMIT, offset: configs.length, search: searchTerm, append: true })
  }

  const handleRowClick = async (config: Config) => {
    await getConfigById(config.id)
    setOpenViewDialog(true)
  }

  const handleClearFilters = () => {
    setSearchInput('')
  }

  const handleCreate = () => {
    setEditData(null)
    setOpenFormDialog(true)
  }

  const handleEdit = async (config: Config) => {
    const fullConfig = await getConfigById(config.id)
    if (fullConfig) {
      setEditData(fullConfig)
      setOpenFormDialog(true)
    }
  }

  const handleDeleteClick = (config: Config) => {
    setDeleteId(config.id)
    setOpenConfirmDialog(true)
  }

  const handleConfirmDelete = async () => {
    if (deleteId) {
      const success = await deleteConfig(deleteId)
      if (success) {
        showNotification('Config deleted successfully', 'success')
      }
      setOpenConfirmDialog(false)
      setDeleteId('')
      loadConfigs({ limit: LIMIT, offset: 0, search: searchTerm, append: false })
    }
  }

  const handleSave = async (data: CreateConfigRequest | UpdateConfigRequest, isEdit: boolean) => {
    if (isEdit && editData) {
      await updateConfig(editData.id, data)
    } else {
      await createConfig(data as CreateConfigRequest)
    }
    setEditData(null)
    setOpenFormDialog(false)
    loadConfigs({ limit: LIMIT, offset: 0, search: searchTerm, append: false })
  }

  const handleCloseFormDialog = () => {
    setOpenFormDialog(false)
    setEditData(null)
  }

  const handleCloseViewDialog = () => {
    setOpenViewDialog(false)
    clearSelectedConfig()
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <PageHeader
        title="Configs"
        subtitle="Manage your API gateway route configurations"
        action={
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={handleCreate}
          >
            Add Config
          </Button>
        }
      />

      {/* Search */}
      <Box sx={{ mb: 3 }}>
        <TextField
          placeholder="Search configs..."
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          size="small"
          sx={{ width: 300 }}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon color="action" />
              </InputAdornment>
            ),
          }}
        />
        {searchTerm && (
          <Button onClick={handleClearFilters} sx={{ ml: 1 }}>
            Clear
          </Button>
        )}
      </Box>

      {/* Table */}
      <ConfigTable
        configs={configs}
        loading={loading}
        loadingMore={loadingMore}
        hasMore={hasMore}
        onRowClick={handleRowClick}
        onEdit={handleEdit}
        onDelete={handleDeleteClick}
        onLoadMore={handleLoadMore}
      />

      {/* Form Dialog */}
      <ConfigFormDialog
        open={openFormDialog}
        onClose={handleCloseFormDialog}
        onSave={handleSave}
        editData={editData}
      />

      {/* View Dialog */}
      <ConfigViewDialog
        open={openViewDialog}
        onClose={handleCloseViewDialog}
        config={selectedConfig}
        onEdit={() => {
          handleCloseViewDialog()
          if (selectedConfig) handleEdit(selectedConfig)
        }}
      />

      {/* Confirm Delete Dialog */}
      <ConfirmDialog
        open={openConfirmDialog}
        title="Delete Config"
        message="Are you sure you want to delete this config? This action cannot be undone."
        onConfirm={handleConfirmDelete}
        onCancel={() => {
          setOpenConfirmDialog(false)
          setDeleteId('')
        }}
      />
    </Container>
  )
}
