import { Route, Routes } from "react-router-dom";
import AuthForm from "./components/AuthForm";

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<>Home</>}></Route>
      <Route path="/login" element={<AuthForm mode="login" />}></Route>
      <Route path="/register" element={<AuthForm mode="register" />}></Route>
    </Routes>
  );
}
