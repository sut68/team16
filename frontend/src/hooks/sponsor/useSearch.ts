import { ref, computed, watch, type Ref } from 'vue'

interface UseSearchOptions {
  debounceMs?: number
}

export function useSearch<T>(
  items: Ref<T[]>,
  filterFn: (item: T, term: string) => boolean,
  options: UseSearchOptions = {}
) {
  const { debounceMs = 300 } = options

  const searchQuery = ref('')
  const q = ref('')

  let timer: ReturnType<typeof setTimeout> | null = null

  // Debounced search
  watch(searchQuery, (val) => {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      q.value = val.trim().toLowerCase()
      timer = null
    }, debounceMs)
  })

  // Filtered items (sorted by ID ASC - 1, 2, 3...)
  const filtered = computed(() => {
    const term = q.value
    let result = term
      ? items.value.filter(item => filterFn(item, term))
      : [...items.value]

    // Sort by ID ASC (1, 2, 3...)
    return result.sort((a: any, b: any) => (a.ID || 0) - (b.ID || 0))
  })

  function clearSearch() {
    searchQuery.value = ''
    q.value = ''
  }

  return {
    searchQuery,
    q,
    filtered,
    clearSearch,
  }
}