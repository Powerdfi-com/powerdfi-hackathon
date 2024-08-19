"use client"
import Loading from '@/app/components/loading';
import AdminAPI from '@/app/utils/apis/admin';
import React, { useEffect, useState } from 'react'
import { toast } from 'react-toastify';

const Notifications = () => {
    const { data } = AdminAPI.getNotificationsPrefs();
    const [prefs, setPrefs] = useState<{ created: boolean, login: boolean }>();
    useEffect(() => {
        if (data) {
            setPrefs(data.data);
        }
    }, [data]);

    const { mutateAsync } = AdminAPI.updateAdminPrefs();

    const handleClickUpdatePrefs = async (data: any) => {
        await toast.promise(mutateAsync(data), {
            pending: "Updating Preference",
            error: "Something went wrong...",
            success: "Prefs updated!"
        }).then((res) => {
            setPrefs({ ...prefs, ...data })
        })
    }
    return <section>
        <h3 className='text-2xl text-white leading-relaxed'>Notification settings</h3>

        {
            prefs ? <div className='w-1/2 py-6'>
                <div className="flex items-center p-2 px-4 border border-border rounded-sm bg-transparent">
                    <label className="w-full py-2 text-sm font-bold text-white/80">Login</label>
                    <input onChange={(e) => handleClickUpdatePrefs({ created: e.target.checked })} checked={prefs.created} type="checkbox" className="w-5 h-5 text-white bg-secondary border-secondary rounded-sm" />
                </div>
                <div className="flex items-center p-2 px-4 border border-border rounded-sm bg-transparent">
                    <label className="w-full py-2 text-sm font-bold text-white/80">Asset Created</label>
                    <input onChange={(e) => handleClickUpdatePrefs({ login: e.target.checked })} checked={prefs.login} type="checkbox" className="w-5 h-5 text-white bg-secondary border-secondary rounded-sm" />
                </div>

            </div> : <Loading />
        }

    </section>
}

export default Notifications