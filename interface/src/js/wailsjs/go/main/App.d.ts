// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {model} from '../models';

export function AnswerChat(arg1:string,arg2:string,arg3:Array<number>):Promise<model.AnswerChatResponse>;

export function AreKeyExist():Promise<boolean>;

export function ConfirmStartOver():Promise<string>;

export function EndChat(arg1:string,arg2:string):Promise<model.AnswerChatResponse>;

export function StartChat(arg1:string,arg2:Array<string>,arg3:string):Promise<model.StartChatResponse>;

export function Status():Promise<model.StatusResponse>;

export function UpdateAPIKeys(arg1:string,arg2:string):Promise<void>;
