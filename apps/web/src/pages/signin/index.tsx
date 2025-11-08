import { setLocalStorageItem } from "@/lib/constants/localStorage";
import { env } from "@/config/env";
import { LOCAL_STORAGE_KEYS } from "../../lib/constants/localStorage";
import { SigninPage } from "@/components/pages/Auth/SigninPage";
import { useSignInPageRedirect } from "@/components/pages/Auth/useSignInPageRedirect";

const SignInPage = () => {
    const { isLoading } = useSignInPageRedirect();

    const handleRequestGoogleOAuthUrl = async () => {
        if (isLoading) {
            return;
        }
        // open new tab with the url
        setLocalStorageItem(
            LOCAL_STORAGE_KEYS.ON_GOOGLE_OAUTH_SUCCESS_REDIRECT,
            `/signin/verify-oauth`,
        );
        window.location.href = `${env.VITE_CORE_BACKEND_API}/api/v1/auth/request-google-oauth`;
    };

    return <SigninPage onGoogleOAuthClick={handleRequestGoogleOAuthUrl} isLoading={isLoading} />;
};

export default SignInPage;
