import { env } from "@/config/env";
import { LOCAL_STORAGE_KEYS, setLocalStorageItem } from "@/lib/constants/localStorage";
import { OnboardMethods } from "../onboard/[method]";
import { SignupPage } from "@/components/pages/Auth/SignupPage";
import { useAuthRedirect } from "@/components/pages/Auth/useAuthRedirect";

const SignUpPage = () => {

    useAuthRedirect();
    const handleRequestGoogleOAuthUrl = async () => {
        // open new tab with the url
        setLocalStorageItem(
            LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT,
            `/onboard/${OnboardMethods.GOOGLE}`,
        );
        window.location.href = `${env.VITE_CORE_BACKEND_API}/api/v1/auth/request-google-oauth`;
    };

    useEffect(() => {
        if (user) {
            navigate("/app");
        }
    }, [user, navigate]);

    return (
        <SignupPage onGoogleOAuthClick={handleRequestGoogleOAuthUrl} />
    )
};

export default SignUpPage;
