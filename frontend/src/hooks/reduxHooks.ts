// https://redux.js.org/tutorials/typescript-quick-start
// Read Define Typed Hooks

import { useDispatch, useSelector } from 'react-redux'
import type { AppDispatch, RootState } from '../store/index'

// Use throughout your app instead of plain `useDispatch` and `useSelector`
export const useAppDispatch = useDispatch.withTypes<AppDispatch>()
export const useAppSelector = useSelector.withTypes<RootState>()