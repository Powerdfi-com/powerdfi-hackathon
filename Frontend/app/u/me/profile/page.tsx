"use client"
import FileUpload from '@/app/components/fileUpload'
import UserAPI from '@/app/utils/apis/user'
import React, { useContext } from 'react'
import { useState, useEffect } from "react"
import { toast } from 'react-toastify'
import { uploadFile } from '@/app/utils/apis/file-upload'
import Image from 'next/image'
import getFileUrl from '@/app/utils/fileUrl'
import { UserContext } from '@/app/context/context'
import { useRouter } from 'next/navigation'


const Profile = () => {
    const [file, setFile] = useState<FileList>();
    const { user, setUser } = useContext(UserContext);
    const router = useRouter()
    const [userData, setUserData] = useState(user);
    const { mutateAsync: updateProfile, isPending: isUpdating } = UserAPI.updateProfile({
        onSuccess: (res) => {
            setUser(res.data);
            setIsEditing(false);
        }
    });

    const handleClickUpdateProfile = async (e: any) => {
        e.preventDefault();
        const info = {
            pending: "Updating profile, please wait!",
            error: "Something went wrong, please try again!",
            success: "Profile successfully updated!"
        }
        const profilePromise = new Promise(async (resolve, reject) => {
            await uploadFile(file!).then(async (value) => {
                const url = value.at(0);
                await updateProfile({ ...userData, avatar: url! }).then(() => {
                    resolve(true);
                }).catch((e) => {
                    reject()
                })
            }).catch(() => {
                reject();
            });
        });
        if (file) {
            await toast.promise(profilePromise, info)
        }
        else {
            await toast.promise(updateProfile(userData), info)
        }
    }

    const [isEditing, setIsEditing] = useState(false);

    return <section className="h-full">
        <h3 className='text-2xl text-white leading-relaxed'>Profile Settings</h3>
        <div className='flex h-full gap-8 mt-8'>
            <div className='w-64 aspect-square'>
                <div className='ring-1 ring-border rounded-lg h-64 w-full flex items-center justify-center relative'>
                    {
                        (user.avatar || file) && <div className='absolute top-2 left-2 right-2 bottom-2'>
                            <div className='relative w-full h-full'>
                                <Image src={file ? getFileUrl(file.item(0)!) : user.avatar} alt="user" fill={true} className='rounded-md object-cover' />
                            </div>
                        </div>
                    }
                    <FileUpload updateValue={(e) => setFile(e)}>
                        <div className='p-2 rounded-lg ring-1 ring-primary/40 text-sm text-white'>Choose File</div>
                    </FileUpload>
                </div>
            </div>
            <form className='flex-1 flex flex-col gap-2'>
                <label className='flex flex-col gap-2'>
                    <span className='text-sm text-white'>Username</span>
                    <input disabled={!isEditing} className='max-w-sm outline-none text-sm text-white border rounded-md border-text-grey bg-transparent py-2 px-5' value={userData.username} onChange={(e) => setUserData({ ...userData, username: e.target.value })} />
                </label>
                <label className='flex flex-col gap-2'>
                    <span className='text-sm text-white'>Bio</span>
                    <textarea disabled={!isEditing} className='max-w-sm outline-none text-sm text-white border rounded-md border-text-grey bg-transparent py-2 px-5' value={userData.bio} onChange={(e) => setUserData({ ...userData, bio: e.target.value })} />
                </label>
                <label className='flex flex-col gap-2'>
                    <span className='text-sm text-white'>Email</span>
                    <input type='email' className='max-w-sm outline-none text-sm text-white border rounded-md border-text-grey bg-transparent py-2 px-5' value={user.email} disabled />
                </label>
                <h4 className='text-white text-lg leading-relaxed mt-4'>Social Profile</h4>
                <label className='flex flex-col gap-2'>
                    <span className='text-sm text-white'>Twitter</span>
                    <input disabled={!isEditing} className='max-w-sm outline-none text-sm text-white border rounded-md border-text-grey bg-transparent py-2 px-5' value={userData.twitter} onChange={(e) => setUserData({ ...userData, twitter: e.target.value })} />
                </label>
                <label className='flex flex-col gap-2'>
                    <span className='text-sm text-white'>Discord</span>
                    <input disabled={!isEditing} className='max-w-sm outline-none text-sm text-white border rounded-md border-text-grey bg-transparent py-2 px-5' value={userData.discord} onChange={(e) => setUserData({ ...userData, discord: e.target.value })} />
                </label>
                <label className='flex flex-col gap-2'>
                    <span className='text-sm text-white'>Website Url</span>
                    <input disabled={!isEditing} className='max-w-sm outline-none text-sm text-white border rounded-md border-text-grey bg-transparent py-2 px-5' value={userData.website} onChange={(e) => setUserData({ ...userData, website: e.target.value })} />
                </label>
                <button disabled={isUpdating} className="text-sm bg-secondary relative py-3 rounded-md w-64 my-6" onClick={isEditing ? handleClickUpdateProfile : (e) => { e.preventDefault(); setIsEditing(true) }}>{isEditing ? "Save" : "Edit"}</button>
            </form>
        </div>
    </section>
}

export default Profile