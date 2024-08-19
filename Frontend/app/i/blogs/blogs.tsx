"use client"
import { Link } from '@/app/components/link'
import Loading from '@/app/components/loading'
import { ghost, GhostAPI } from '@/app/utils/apis/ghost'
import { dateFromISO, formatDate } from '@/app/utils/func'
import { useMutation } from '@tanstack/react-query'
import { PostsOrPages } from '@tryghost/content-api'
import Image from 'next/image'
import React, { useEffect, useState } from 'react'
import { IoPersonSharp } from 'react-icons/io5'


const GhostBlogs = () => {
    const [posts, setPosts] = useState<PostsOrPages>();
    const { data, mutateAsync } = GhostAPI.fetchBlogs();
    useEffect(() => {
        const fetchData = async () => {
            await mutateAsync().then((res) => {
                setPosts(res)
            })
        }
        fetchData()
    }, [mutateAsync])
    if (!data) {
        return <Loading />
    }
    return <section className='px-24 py-12'>
        <section className="relative flex flex-col justify-center items-center py-12 w-full">
            <div className="relative w-full h-[560px] flex justify-center items-center py-8 rounded-2xl">
                <Image src={data![0].feature_image!} fill alt="logo" className="flex justify-center items-center w-full rounded-2xl" />
            </div>
            <div className="absolute  left-24 -bottom-24 bg-background rounded-xl text-white text-[36px] leading-tight w-[598px] h-[260px] font-semibold px-12 justify-center flex flex-col items-start">
                <h3 className="flex justify-center items-center">
                    <q>{data![0].title}</q>
                </h3>
                <div className="flex justify-start items-center text-[14px] text-text-ash space-x-8 py-4">
                    {
                        data[0].authors && <div className="flex justify-center items-center space-x-3">
                            <div className="rounded-full bg-white/40 p-2">
                                <IoPersonSharp className="!text-xl text-black" />
                            </div>
                            <p className="">{data[0].authors[0].name}</p>
                        </div>
                    }
                    <p>{formatDate(dateFromISO(data![0].updated_at!))}</p>
                </div>
            </div>
        </section>
        <div className='py-2 px-4 text-[40px] rounded-md ring-1 ring-primary/40 text-secondary font-semibold gradient w-fit mt-48'>Latest Posts</div>
        {
            posts ? <ul className='grid grid-cols-3 mt-12'>
                {posts.map((post) => <Link href={`/i/blogs/${post.slug}`} key={post.id} className='w-[330px] h-[320px] bg-border rounded-lg'>
                    <div className='w-full relative h-[257px] '>
                        <Image fill src={post.feature_image!} alt={post.slug} className='object-cover rounded-t-lg' />
                    </div>
                    <div className='h-[63px] w-full flex px-6 items-center text-left text-white text-[20px]'>{post.title}</div>
                </Link>)}
            </ul> : <Loading />
        }
    </section>
}

export default GhostBlogs