import { CustomHelmet } from '../providers/helmets/CustomHelmet';
import { useHelmet } from '../../hooks/useHelmet';

export const LoginPage = () => {
    const helmetData = useHelmet({
        pageType: 'login',
        title: 'Sign In | DECM Platform',
        description: 'Sign in to access your DECM account and manage your digital credentials'
    });

    return (
        <>
            <CustomHelmet
                title={helmetData.title}
                description={helmetData.description}
                themeColor={helmetData.themeColor}
                additionalMeta={[
                    { name: 'robots', content: 'noindex, nofollow' }, // Don't index login pages
                    { property: 'og:type', content: 'website' },
                    { property: 'og:title', content: helmetData.title },
                    { property: 'og:description', content: helmetData.description }
                ]}
            />
            <div>
                <h1>Login to DECM</h1>
                <p>Sign in to your account</p>
                {/* Your login form would go here */}
            </div>
        </>
    );
};