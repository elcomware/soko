import { ref, computed } from 'vue'

export interface Currency {
  code: string
  symbol: string
  name: string
  exchange_rate: number
  prefix: boolean
}

export interface MarketRegion {
  id: string
  name: string
  flag: string
  default_currency: string
  local_logistics: string
  payment_rails: string[]
  vat_percentage: number
  customs_clearance_hr: number
}

// Global reactive currency state
const selectedCurrency = ref<string>('USD')
const selectedRegion = ref<string>('global')

export const defaultCurrencies: Currency[] = [
  { code: 'USD', symbol: '$', name: 'US Dollar', exchange_rate: 1.0, prefix: true },
  { code: 'EUR', symbol: '€', name: 'Euro', exchange_rate: 0.92, prefix: true },
  { code: 'GBP', symbol: '£', name: 'British Pound', exchange_rate: 0.78, prefix: true },
  { code: 'KES', symbol: 'KSh', name: 'Kenya Shilling', exchange_rate: 130.0, prefix: false },
  { code: 'NGN', symbol: '₦', name: 'Nigerian Naira', exchange_rate: 1550.0, prefix: true },
  { code: 'ZAR', symbol: 'R', name: 'South African Rand', exchange_rate: 18.5, prefix: true },
  { code: 'AED', symbol: 'AED', name: 'UAE Dirham', exchange_rate: 3.67, prefix: false },
  { code: 'CNY', symbol: '¥', name: 'Chinese Yuan', exchange_rate: 7.23, prefix: true },
  { code: 'JPY', symbol: '¥', name: 'Japanese Yen', exchange_rate: 152.0, prefix: true },
  { code: 'INR', symbol: '₹', name: 'Indian Rupee', exchange_rate: 84.0, prefix: true },
]

export function useCurrency(initialCurrencies?: Currency[]) {
  const currencies = ref<Currency[]>(initialCurrencies && initialCurrencies.length ? initialCurrencies : defaultCurrencies)

  const activeCurrency = computed<Currency>(() => {
    return currencies.value.find(c => c.code === selectedCurrency.value) || currencies.value[0]
  })

  const setCurrency = (code: string) => {
    selectedCurrency.value = code
  }

  const setRegion = (regionId: string, defaultCurr?: string) => {
    selectedRegion.value = regionId
    if (defaultCurr) {
      selectedCurrency.value = defaultCurr
    }
  }

  const formatPrice = (usdAmount: number): string => {
    if (usdAmount === undefined || usdAmount === null) return '$0.00'
    const curr = activeCurrency.value
    const converted = usdAmount * curr.exchange_rate

    const formattedNumber = converted.toLocaleString(undefined, {
      minimumFractionDigits: converted >= 1000 ? 0 : 2,
      maximumFractionDigits: 2,
    })

    if (curr.prefix) {
      return `${curr.symbol}${formattedNumber}`
    }
    return `${formattedNumber} ${curr.symbol}`
  }

  return {
    selectedCurrency,
    selectedRegion,
    currencies,
    activeCurrency,
    setCurrency,
    setRegion,
    formatPrice,
  }
}
