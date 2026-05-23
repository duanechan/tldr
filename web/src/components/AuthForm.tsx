import {
    EyeClosedIcon,
    EyeIcon,
    SignInIcon,
    UserPlusIcon,
} from "@phosphor-icons/react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "./ui/button";
import {
    Card,
    CardAction,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "./ui/card";
import { Input } from "./ui/input";
import { InputGroup, InputGroupAddon, InputGroupInput } from "./ui/input-group";
import { Label } from "./ui/label";
import { ErrorResponseSchema } from "@/lib/schema";
import { TokenKey } from "@/lib/constants";

export default function AuthForm({ mode }: { mode: "login" | "register" }) {
    const navigate = useNavigate();
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [isPasswordShown, setPasswordShown] = useState(false);
    const [isConfirmPasswordShown, setConfirmPasswordShown] = useState(false);
    const [errors, setErrors] = useState<string[]>([]);

    async function handleLogin(e: React.MouseEvent) {
        e.preventDefault();
        setErrors([]);

        try {
            const res = await fetch("/api/v1/auth/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password }),
            });
            const data = await res.json();
            if (ErrorResponseSchema.safeParse(data).success) {
                setErrors((prev) => [...prev, data.message]);
                return;
            }

            localStorage.setItem(TokenKey, data);
            navigate("/");
        } catch (e) {
            if (e instanceof Error) {
                setErrors((prev) => [...prev, e.message]);
            }
        }
    }

    async function handleRegister(e: React.MouseEvent) {
        e.preventDefault();
        setErrors([]);

        try {
            const res = await fetch("/api/v1/auth/register", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password, confirmPassword }),
            });
            const data = await res.json();
            if (ErrorResponseSchema.safeParse(data).success) {
                setErrors((prev) => [...prev, data.message]);
                return;
            }

            localStorage.setItem(TokenKey, data);
            navigate("/");
        } catch (e) {
            if (e instanceof Error) {
                setErrors((prev) => [...prev, e.message]);
            }
        }
    }

    return (
        <Card className="w-full max-w-md border rounded-xl border-mist-700">
            <CardHeader>
                <CardTitle className="text-lg font-bold">
                    {mode === "login" ? "Login" : "Register"}
                </CardTitle>
                <CardDescription>
                    {mode === "login"
                        ? "Sign in to your account."
                        : "Create an account."}
                </CardDescription>
                <CardAction
                    onClick={() => {
                        navigate(mode === "login" ? "/register" : "/login");
                    }}
                >
                    <Button type="button" variant="link" aria-label="Sign Up">
                        {mode === "login" ? (
                            <>
                                Sign Up
                                <UserPlusIcon className="size-5" />
                            </>
                        ) : (
                            <>
                                Sign In
                                <SignInIcon className="size-6" />
                            </>
                        )}
                    </Button>
                </CardAction>
            </CardHeader>
            <form>
                <CardContent>
                    <div className="flex flex-col gap-4">
                        <div className="flex flex-col gap-2">
                            <Label htmlFor="usernameField">Username</Label>
                            <Input
                                id="usernameField"
                                type="text"
                                placeholder="johndoe67"
                                onChange={(e) => setUsername(e.target.value)}
                                className="rounded-lg border-mist-700"
                                required
                            />
                        </div>
                        <div className="flex flex-col gap-2">
                            <Label htmlFor="passwordField">Password</Label>
                            <InputGroup>
                                <InputGroupInput
                                    id="passwordField"
                                    type={isPasswordShown ? "text" : "password"}
                                    placeholder="Enter your password"
                                    onChange={(e) =>
                                        setPassword(e.target.value)
                                    }
                                    required
                                />
                                <InputGroupAddon
                                    align="inline-end"
                                    className="cursor-pointer"
                                    onClick={() =>
                                        setPasswordShown(!isPasswordShown)
                                    }
                                >
                                    {isPasswordShown ? (
                                        <EyeIcon />
                                    ) : (
                                        <EyeClosedIcon />
                                    )}
                                </InputGroupAddon>
                            </InputGroup>
                        </div>
                        <div
                            className={`overflow-hidden transition-all duration-300 ${mode === "register" ? "max-h-24" : "max-h-0"}`}
                        >
                            <div className="flex flex-col gap-2">
                                <Label htmlFor="confirmPasswordField">
                                    Confirm Password
                                </Label>
                                <InputGroup>
                                    <InputGroupInput
                                        id="confirmPasswordField"
                                        type={
                                            isConfirmPasswordShown
                                                ? "text"
                                                : "password"
                                        }
                                        placeholder="Re-enter your password"
                                        onChange={(e) =>
                                            setConfirmPassword(e.target.value)
                                        }
                                        required
                                    />
                                    <InputGroupAddon
                                        align="inline-end"
                                        className="cursor-pointer"
                                        onClick={() =>
                                            setConfirmPasswordShown(
                                                !isConfirmPasswordShown,
                                            )
                                        }
                                    >
                                        {isConfirmPasswordShown ? (
                                            <EyeIcon />
                                        ) : (
                                            <EyeClosedIcon />
                                        )}
                                    </InputGroupAddon>
                                </InputGroup>
                            </div>
                        </div>
                    </div>
                </CardContent>
                <CardFooter className="flex-col pt-6">
                    <Button
                        type="submit"
                        className="w-full rounded-lg duration-300 transition-all hover:-translate-y-0.5"
                        onClick={(e) =>
                            mode === "login"
                                ? handleLogin(e)
                                : handleRegister(e)
                        }
                    >
                        {mode === "login" ? "Login" : "Register"}
                    </Button>
                    {errors.length !== 0 && errors.map((e) => <span>{e}</span>)}
                </CardFooter>
            </form>
        </Card>
    );
}
