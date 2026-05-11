import { Route, Routes } from "react-router-dom";
import AuthForm from "./components/AuthForm";
import AppLayout from "./layouts/AppLayout";
import ProtectedLayout from "./layouts/ProtectedLayout";

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<AuthForm mode="login" />} />
      <Route path="/register" element={<AuthForm mode="register" />} />

      <Route element={<ProtectedLayout />}>
        <Route element={<AppLayout />}>
          <Route path="/" element={<>Home</>} />
        </Route>
      </Route>
    </Routes>
  );
}
