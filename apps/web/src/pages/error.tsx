import { FaviconHelmet } from '@/components/providers/helmets/FaviconHelmet';
import { Error } from '@/components/pages/Error';
import { useTranslation } from 'react-i18next';

const ErrorPage = () => {
    const { t } = useTranslation();

    return (
        <>
            <FaviconHelmet
                title={`${t('error.title')} | ${t('common.appName')}`}
                description={t('error.description')}
            />
            <Error />
        </>
    );
};

export default ErrorPage;

