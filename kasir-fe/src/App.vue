<script setup>
import { ref, onMounted } from 'vue'
import {
  LayoutDashboard,
  Package,
  Users,
  Layers,
  History,
  PlusCircle,
  AlertTriangle,
  Sprout,
  Menu
} from 'lucide-vue-next'

import DashboardView from './components/DashboardView.vue'
import MasterBarangView from './components/MasterBarangView.vue'
import MasterKontakView from './components/MasterKontakView.vue'
import StokBatchView from './components/StokBatchView.vue'
import TrackingMutasiView from './components/TrackingMutasiView.vue'
import TransaksiModal from './components/TransaksiModal.vue'
import { reportApi } from './services/api'

const currentView = ref('dashboard')
const isMutasiModalOpen = ref(false)
const selectedBarangIdForMutasi = ref(null)
const lowStockCount = ref(0)
const isMobileMenuOpen = ref(false)

const dashboardRef = ref(null)
const barangRef = ref(null)
const batchRef = ref(null)
const mutasiRef = ref(null)

const navItems = [
  { id: 'dashboard', label: 'Dashboard Overview', icon: LayoutDashboard },
  { id: 'barang', label: 'Master Produk', icon: Package },
  { id: 'kontak', label: 'Master Kontak (Party)', icon: Users },
  { id: 'batch', label: 'Laporan Batch FEFO', icon: Layers },
  { id: 'tracking', label: 'Audit Mutasi Stok', icon: History }
]

const fetchLowStockCount = async () => {
  try {
    const res = await reportApi.getLowStockAlerts()
    lowStockCount.value = res.data.total_alert || 0
  } catch (err) {
    console.error(err)
  }
}

const openMutasi = (barangId = null) => {
  selectedBarangIdForMutasi.value = barangId
  isMutasiModalOpen.value = true
}

const handleMutasiSuccess = () => {
  fetchLowStockCount()
  dashboardRef.value?.refresh?.()
  barangRef.value?.refresh?.()
  batchRef.value?.refresh?.()
  mutasiRef.value?.refresh?.()
}

const handleNavigate = (viewId) => {
  currentView.value = viewId
  isMobileMenuOpen.value = false
}

onMounted(() => {
  fetchLowStockCount()
})
</script>

<template>
  <div class="min-h-screen flex bg-slate-50 text-slate-900">
    <!-- Sidebar (240px Uncodixfy Standard) -->
    <aside
      :class="[
        'fixed inset-y-0 left-0 z-40 w-60 bg-white border-r border-slate-200 flex flex-col transition-transform duration-200 lg:translate-x-0',
        isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full'
      ]"
    >
      <!-- Brand Logo -->
      <div class="h-16 flex items-center gap-2.5 px-5 border-b border-slate-100">
        <div class="w-8 h-8 rounded-lg bg-emerald-600 flex items-center justify-center text-white shadow-sm">
          <Sprout class="w-5 h-5" />
        </div>
        <div>
          <h1 class="text-sm font-bold text-slate-900 leading-tight">AgriStock</h1>
          <p class="text-[10px] text-slate-500 font-medium">Inventaris Toko Pertanian</p>
        </div>
      </div>

      <!-- Navigation Links -->
      <nav class="flex-1 px-3 py-4 space-y-1">
        <button
          v-for="item in navItems"
          :key="item.id"
          @click="handleNavigate(item.id)"
          :class="[
            'w-full flex items-center gap-3 px-3 py-2 rounded-lg text-xs font-medium transition',
            currentView === item.id
              ? 'bg-emerald-50 text-emerald-700 font-semibold'
              : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
          ]"
        >
          <component :is="item.icon" class="w-4 h-4 shrink-0" />
          <span>{{ item.label }}</span>
          <span
            v-if="item.id === 'dashboard' && lowStockCount > 0"
            class="ml-auto px-1.5 py-0.5 text-[10px] font-bold bg-amber-100 text-amber-800 rounded"
          >
            {{ lowStockCount }}
          </span>
        </button>
      </nav>

      <!-- Quick Action in Sidebar Bottom -->
      <div class="p-3 border-t border-slate-100">
        <button
          @click="openMutasi()"
          class="w-full flex items-center justify-center gap-2 py-2 px-3 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-xs font-medium shadow-sm transition"
        >
          <PlusCircle class="w-4 h-4" />
          <span>Transaksi Mutasi</span>
        </button>
      </div>
    </aside>

    <!-- Mobile Overlay Backdrop -->
    <div
      v-if="isMobileMenuOpen"
      @click="isMobileMenuOpen = false"
      class="fixed inset-0 z-30 bg-slate-900/50 lg:hidden"
    ></div>

    <!-- Main Content Area -->
    <div class="flex-1 lg:pl-60 flex flex-col min-w-0">
      <!-- Top Navbar -->
      <header class="h-16 bg-white border-b border-slate-200 sticky top-0 z-20 flex items-center justify-between px-4 sm:px-8">
        <div class="flex items-center gap-3">
          <button
            @click="isMobileMenuOpen = !isMobileMenuOpen"
            class="p-1.5 rounded-lg text-slate-600 hover:bg-slate-100 lg:hidden"
          >
            <Menu class="w-5 h-5" />
          </button>
          <div>
            <h2 class="text-sm font-semibold text-slate-900 capitalize">
              {{ navItems.find(n => n.id === currentView)?.label }}
            </h2>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <!-- Low Stock Alert Header Badge -->
          <button
            @click="handleNavigate('dashboard')"
            v-if="lowStockCount > 0"
            class="inline-flex items-center gap-1.5 px-2.5 py-1 bg-amber-50 border border-amber-200 text-amber-800 rounded-lg text-xs font-medium hover:bg-amber-100 transition"
          >
            <AlertTriangle class="w-3.5 h-3.5 text-amber-600" />
            <span>{{ lowStockCount }} Stok Menipis</span>
          </button>

          <!-- Top Mutation Trigger -->
          <button
            @click="openMutasi()"
            class="hidden sm:inline-flex items-center gap-1.5 px-3 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-xs font-medium shadow-sm transition"
          >
            <PlusCircle class="w-3.5 h-3.5" />
            <span>+ Mutasi Stok</span>
          </button>
        </div>
      </header>

      <!-- View Render Area -->
      <main class="flex-1 p-4 sm:p-8 max-w-7xl w-full mx-auto">
        <DashboardView
          v-if="currentView === 'dashboard'"
          ref="dashboardRef"
          @openMutasi="openMutasi"
          @navigate="handleNavigate"
        />

        <MasterBarangView
          v-else-if="currentView === 'barang'"
          ref="barangRef"
          @openMutasiFor="openMutasi"
        />

        <MasterKontakView
          v-else-if="currentView === 'kontak'"
        />

        <StokBatchView
          v-else-if="currentView === 'batch'"
          ref="batchRef"
        />

        <TrackingMutasiView
          v-else-if="currentView === 'tracking'"
          ref="mutasiRef"
        />
      </main>
    </div>

    <!-- Global Stock Mutation Modal -->
    <TransaksiModal
      :isOpen="isMutasiModalOpen"
      :preselectedBarangId="selectedBarangIdForMutasi"
      @close="isMutasiModalOpen = false"
      @success="handleMutasiSuccess"
    />
  </div>
</template>
