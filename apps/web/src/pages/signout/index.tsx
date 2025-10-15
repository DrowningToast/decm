import { useSignout } from "@/components/useSignout";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

const SignoutPage = () => {

    const { signout, isPending, error } = useSignout();
    const navigate = useNavigate();
    useEffect(() => {
        const init = async () => {
            await signout();
            navigate("/");
        }
        init();
    }, [navigate, signout]);

    if (isPending) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error.message}</div>;
    }

    return (
        <div>
            <h1>You're logging out. Please wait...</h1>
        </div>
    )
}

export default SignoutPage;