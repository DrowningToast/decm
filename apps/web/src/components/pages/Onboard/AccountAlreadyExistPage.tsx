import { Link } from "@/router";

export const AccountAlreadyExistsPage = () => {

    return (
        <>
            <h1>[PH] Account Already Exists Page. Please sign in </h1>
            <Link to="/signin">
                Go to sign in
            </Link>
        </>
    )
}
