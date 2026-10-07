export type User={id:number;name:string;email:string};
export type Book={id:number;title:string;author:string;description:string;cover_url:string;genre:string;average_rating:number;review_count:number};
export type Review={id:number;book_id:number;user_id:number;user_name:string;rating:number;body:string;created_at:string};
export type LibraryBook=Book&{status:'want'|'reading'|'read'};
const base=process.env.NEXT_PUBLIC_API_URL||'http://localhost:8080/api';
export const token=()=>typeof window==='undefined'?'':localStorage.getItem('qare_token')||'';
export async function api<T>(path:string,options:RequestInit={}):Promise<T>{const response=await fetch(base+path,{...options,headers:{'Content-Type':'application/json',...(token()?{Authorization:`Bearer ${token()}`}:{ }),...options.headers},cache:'no-store'});if(!response.ok){const body=await response.json().catch(()=>({}));throw new Error(body.error||'حدث خطأ، حاول مرة أخرى')};if(response.status===204)return undefined as T;return response.json()}
export const auth={login:(email:string,password:string)=>api<{user:User;token:string}>('/auth/login',{method:'POST',body:JSON.stringify({email,password})}),register:(name:string,email:string,password:string)=>api<{user:User;token:string}>('/auth/register',{method:'POST',body:JSON.stringify({name,email,password})}),me:()=>api<User>('/auth/me')};
export const books={list:(q='')=>api<Book[]>('/books?q='+encodeURIComponent(q)),get:(id:string)=>api<{book:Book;reviews:Review[]}>('/books/'+id),review:(id:string,rating:number,body:string)=>api('/books/'+id+'/review',{method:'PUT',body:JSON.stringify({rating,body})}),deleteReview:(id:string)=>api('/books/'+id+'/review',{method:'DELETE'})};
export const library={list:()=>api<LibraryBook[]>('/library'),save:(id:number,status:string)=>api('/library/'+id,{method:'PUT',body:JSON.stringify({status})}),remove:(id:number)=>api('/library/'+id,{method:'DELETE'})};
