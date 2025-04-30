import { configureStore, createSlice } from "@reduxjs/toolkit";

const initialState = {
  user: localStorage.getItem("user")
    ? JSON.parse(localStorage.getItem("user"))
    : null,
  darkMode: localStorage.getItem("darkmode")
    ? localStorage.getItem("darkmode")
    : false,
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    signIn: (state, action) => {
      state.user = action.payload;
      localStorage.setItem("user", JSON.stringify(action.payload));
    },
    signOut: (state) => {
      state.user = null;
      localStorage.removeItem("user");
    },
  },
});

const themeSlice = createSlice({
  name: "theme",
  initialState,
  reducers: {
    toggleDarkMode: (state) => {
      state.darkMode = !state.darkMode;
    },
    setDarkMode: (state, action) => {
      state.darkMode = action.payload;
      localStorage.setItem("darkmode", action.payload);
    },
  },
});

const store = configureStore({
  reducer: {
    auth: authSlice.reducer,
    theme: themeSlice.reducer,
  },
});

export const { signIn, signOut } = authSlice.actions;
export const { toggleDarkMode, setDarkMode } = themeSlice.actions;
export default store;
